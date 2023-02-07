package models

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/goccy/go-json"
	"io"
	"log"
	"strconv"
	"strings"
	"treehole_next/config"
	"treehole_next/utils"
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
			Index  string              `json:"_index"`
			ID     string              `json:"_id"`
			Score  float64             `json:"_score"`
			Source SearchFloorResponse `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type SearchFloorResponse struct {
	Content string `json:"content"`
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
		floorIDs[i], err = strconv.Atoi(hit.ID)
		if err != nil {
			return nil, &utils.HttpError{Code: 500, Message: "error parse floor_id from elasticsearch ID"}
		}
	}
	log.Printf("search response: %d\n", floorIDs)

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

type FloorModel struct {
	ID      int    `json:"-"`
	Content string `json:"content"`
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
		floorModel := FloorModel{Content: floor.Content}
		// data: should not contain \n, because \n is the delimiter of one action
		data, err := json.Marshal(floorModel)
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
	log.Printf("Preparing insert floors %v\n", floorIDs)

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error indexing floors %v: %s", floorIDs, err)
		return
	}
	log.Printf("index floors %v success\n", floorIDs)
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
	log.Printf("Preparing delete floors %v\n", floorIDs)

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error deleting floors %v: %s", floorIDs, err)
		return
	}
	log.Printf("delete floors %v success\n", floorIDs)
}

// FloorIndex insert or replace a document, used when a floor is created or restored
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-index_.html
func FloorIndex(floorID int, content string) {
	if ES == nil {
		return
	}
	floorModel := FloorModel{Content: content}
	data, err := json.Marshal(&floorModel)
	if err != nil {
		log.Printf("floor encode error: floor_id: %v", floorID)
		return
	}

	req := esapi.IndexRequest{
		Index:      IndexName,
		DocumentID: strconv.Itoa(floorID),
		Body:       bytes.NewBuffer(data),
		Refresh:    "false",
	}

	res, err := req.Do(context.Background(), ES)
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		data, _ := io.ReadAll(res.Body)
		log.Printf("error index floor: %d: %s\n", floorID, string(data))
	} else {
		log.Printf("index floor success: %d\n", floorID)
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
		log.Printf("delete floor success: %d\n", floorID)
	}
}
