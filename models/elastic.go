package models

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"

	"treehole_next/config"
	"treehole_next/utils"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/goccy/go-json"
)

var ES *elasticsearch.Client

const IndexName = "floors"

func Init() {
	if config.Config.Mode == "test" || config.Config.Mode == "bench" || config.Config.ElasticsearchUrl == "" {
		return
	}

	// export ELASTICSEARCH_URL environment variable to set the ElasticSearch URL
	// example: http://user:pass@127.0.0.1:9200
	var err error
	ES, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.Config.ElasticsearchUrl},
	})
	if err != nil {
		log.Printf("error creating elasticsearch client: %s", err)
		ES = nil
		return
	}

	res, err := ES.Info()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting elasticsearch response")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.IsError() {
		log.Fatal().Str("status", res.Status()).Msg("error getting elasticsearch response")
	}
	var r Map
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatal().Err(err).Msg("Error parsing the elasticsearch response body")
	}

	// print Client and Server Info
	log.Info().Msgf("elasticsearch Client: %s\n", elasticsearch.Version)
	log.Info().Msgf("elasticsearch Server: %s", r["version"].(map[string]interface{})["number"])
}

type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
		Skipped    int `json:"skipped"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			Index  string     `json:"_index"`
			ID     string     `json:"_id"`
			Score  float64    `json:"_score"`
			Source FloorModel `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type FloorModel struct {
	ID        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
}

func Search(keyword string, size, offset int) (Floors, error) {
	if ES == nil {
		return SearchOld(keyword, size, offset)
	}
	req := esapi.SearchRequest{
		Index: []string{IndexName},
		From:  &offset,
		Size:  &size,
		Query: keyword,
		Sort: []string{
			"_score:desc",
			"updated_at:desc",
		},
	}

	res, err := req.Do(context.Background(), ES)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.IsError() {
		var data []byte
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		} else {
			return nil, &common.HttpError{Code: 502, Message: string(data)}
		}
	}

	var response SearchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// get floors
	floorSize := len(response.Hits.Hits)
	floors := make(Floors, 0, floorSize)
	if floorSize == 0 {
		return floors, nil
	}

	floorIDs := make([]int, floorSize)
	for i, hit := range response.Hits.Hits {
		floorIDs[i] = hit.Source.ID
		if err != nil {
			return nil, &common.HttpError{Code: 500, Message: "error parse floor_id from elasticsearch ID"}
		}
	}
	log.Info().Ints("floor_ids", floorIDs).Msg("search response")

	err = DB.Preload("Mention").Find(&floors, floorIDs).Error
	if err != nil {
		return nil, err
	}

	return utils.OrderInGivenOrder(floors, floorIDs), nil
}

func SearchOld(keyword string, size, offset int) (Floors, error) {
	floors := Floors{}
	result := DB.
		Where("content like ?", "%"+keyword+"%").
		Where("hole_id in (?)", DB.Table("hole").Select("id").Where("hidden = false")).
		Offset(offset).Limit(size).Order("id desc").
		Preload("Mention").Find(&floors)
	return floors, result.Error
}

// BulkInsert run in single goroutine only
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html
func BulkInsert(floors []FloorModel) {
	if ES == nil {
		return
	}
	if len(floors) == 0 {
		return
	}

	var BulkBuffer = bytes.NewBuffer(make([]byte, 0, 1024000)) // 100 KB buffer

	for _, floor := range floors {
		// meta: use index, it will insert or replace a document
		BulkBuffer.WriteString(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, floor.ID, "\n"))

		// data: should not contain \n, because \n is the delimiter of one action
		data, err := json.Marshal(floor)
		if err != nil {
			log.Printf("error failed to marshal floor: %s", err)
			return
		}
		BulkBuffer.Write(data)
		BulkBuffer.WriteByte('\n') // the final line of data must end with a newline character \n
	}

	var floorIDs []int
	for _, floorModel := range floors {
		floorIDs = append(floorIDs, floorModel.ID)
	}
	log.Info().Ints("floor_ids", floorIDs).Msg("Preparing insert floors")

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error indexing floors %v: %s", floorIDs, err)
		return
	}
	log.Info().Ints("floor_ids", floorIDs).Msg("index floors success")
}

// BulkDelete used when a hole becomes hidden and delete all of its floors
func BulkDelete(floorIDs []int) {
	if ES == nil {
		return
	}
	if len(floorIDs) == 0 {
		return
	}

	var BulkBuffer = bytes.NewBuffer(make([]byte, 0, 1024000)) // 100 KB buffer

	for _, floorID := range floorIDs {
		// meta: use index, it will insert or replace a document
		BulkBuffer.WriteString(fmt.Sprintf(`{ "delete" : { "_id" : "%d" } }%s`, floorID, "\n"))
	}
	log.Info().Ints("floor_ids", floorIDs).Msg("Preparing delete floors")

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error deleting floors %v: %s", floorIDs, err)
		return
	}
	log.Info().Ints("floor_ids", floorIDs).Msg("delete floors success")
}

// FloorIndex insert or replace a document, used when a floor is created or restored
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-index_.html
func FloorIndex(floorModel FloorModel) {
	if ES == nil {
		return
	}

	data, err := json.Marshal(&floorModel)
	if err != nil {
		log.Err(err).Int("floor_id", floorModel.ID).Msg("floor encode error")
		return
	}

	req := esapi.IndexRequest{
		Index:      IndexName,
		DocumentID: strconv.Itoa(floorModel.ID),
		Body:       bytes.NewBuffer(data),
		Refresh:    "false",
	}

	res, err := req.Do(context.Background(), ES)
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		response, _ := io.ReadAll(res.Body)
		log.Err(err).
			Str("status", res.Status()).
			Int("floor_id", floorModel.ID).
			Bytes("data", response).
			Msg("error index floor")
	} else {
		log.Info().Int("floor_id", floorModel.ID).Msg("index floor success")
	}
}

// FloorDelete used when a floor is deleted
func FloorDelete(floorID int) {
	if ES == nil {
		return
	}
	res, err := ES.Delete(
		IndexName,
		strconv.Itoa(floorID))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		response, _ := io.ReadAll(res.Body)
		log.Err(err).
			Str("status", res.Status()).
			Int("floor_id", floorID).
			Bytes("data", response).
			Msg("error delete floor")
	} else {
		log.Info().Int("floor_id", floorID).Msg("delete floor success")
	}
}
