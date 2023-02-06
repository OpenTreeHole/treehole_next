package elastic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/goccy/go-json"
	"log"
	"strconv"
	"treehole_next/utils"

	. "treehole_next/models"
)

type FloorModel struct {
	Content string `json:"content"`
}

// BulkInsert run in single goroutine only
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html
func BulkInsert(floors Floors) {
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

	floorIDs := utils.Models2IDSlice(floors)
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
func BulkDelete(floors Floors) {
	if len(floors) == 0 {
		return
	}

	var BulkBuffer = bytes.NewBuffer(make([]byte, 0, 1024000)) // 100 KB buffer

	for _, floor := range floors {
		// meta: use index, it will insert or replace a document
		BulkBuffer.WriteString(fmt.Sprintf(`{ "delete" : { "_id" : "%d" } }%s`, floor.ID, "\n"))
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

	floorIDs := utils.Models2IDSlice(floors)
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
	var buffer = bytes.NewBuffer(make([]byte, 16384))

	floorModel := FloorModel{Content: content}
	err := json.NewEncoder(buffer).Encode(floorModel)
	if err != nil {
		log.Printf("floor encode error: floor_id: %v", floorID)
		return
	}

	req := esapi.IndexRequest{
		Index:      IndexName,
		DocumentID: strconv.Itoa(floorID),
		Body:       buffer,
		Refresh:    "false",
	}

	res, err := req.Do(context.Background(), ES)
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error index floor: %d\n", floorID)
	} else {
		log.Printf("index floor success: %d\n", floorID)
	}
}

// FloorDelete used when a floor is deleted
func FloorDelete(floorID int) {
	res, err := ES.Delete(
		IndexName,
		strconv.Itoa(floorID))
	defer func() {
		_ = res.Body.Close()
	}()
	if err != nil || res.IsError() {
		log.Printf("error index floor: %d\n", floorID)
	} else {
		log.Printf("index floor success: %d\n", floorID)
	}
}
