package models

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
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
		log.Printf("Error creating elasticsearch client: %s", err)
		ES = nil
		return
	}

	res, err := ES.Info()
	if err != nil {
		log.Fatalf("Error getting elasticsearch response: %s", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	var r Map
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the elasticsearch response body: %s", err.Error())
	}

	// print Client and Server Info
	log.Printf("elasticsearch Client: %s\n", elasticsearch.Version)
	log.Printf("elasticsearch Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
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
			return nil, &utils.HttpError{Code: 502, Message: string(data)}
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
			return nil, &utils.HttpError{Code: 500, Message: "error parse floor_id from elasticsearch ID"}
		}
	}
	fmt.Printf("search response: %d\n", floorIDs)

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
	fmt.Printf("Preparing insert floors %v\n", floorIDs)

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error indexing floors %v: %s", floorIDs, err)
		return
	}
	fmt.Printf("index floors %v success\n", floorIDs)
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
	fmt.Printf("Preparing delete floors %v\n", floorIDs)

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error deleting floors %v: %s", floorIDs, err)
		return
	}
	fmt.Printf("delete floors %v success\n", floorIDs)
}

// FloorIndex insert or replace a document, used when a floor is created or restored
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-index_.html
func FloorIndex(floorModel FloorModel) {
	if ES == nil {
		return
	}

	data, err := json.Marshal(&floorModel)
	if err != nil {
		log.Printf("floor encode error: floor_id: %v", floorModel.ID)
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
		data, _ := io.ReadAll(res.Body)
		log.Printf("error index floor: %d: %s\n", floorModel.ID, string(data))
	} else {
		fmt.Printf("index floor success: %d\n", floorModel.ID)
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
		data, _ := io.ReadAll(res.Body)
		log.Printf("error delete floor: %d: %s\n", floorID, string(data))
	} else {
		fmt.Printf("delete floor success: %d\n", floorID)
	}
}
