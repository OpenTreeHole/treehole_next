package elastic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/goccy/go-json"
	"log"
	"strconv"

	. "treehole_next/models"
)

type FloorModel struct {
	Content string `json:"content"`
}

// BulkInsert run in single goroutine only, used when dump floors
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html
func BulkInsert(floors Floors) error {
	if len(floors) == 0 {
		return nil
	}

	var BulkBuffer = bytes.NewBuffer(make([]byte, 0, 1024000)) // 100 KB buffer

	firstFloorID := floors[0].ID
	lastFloorID := floors[len(floors)-1].ID
	for _, floor := range floors {
		// meta: use index, it will insert or replace a document
		BulkBuffer.WriteString(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, floor.ID, "\n"))
		// data: should not contain \n, because \n is the delimiter of one action
		data, err := json.Marshal(floor)
		if err != nil {
			return fmt.Errorf("error failed to marshal floor: %s", err)
		}
		BulkBuffer.Write(data)
		BulkBuffer.WriteByte('\n') // the final line of data must end with a newline character \n
	}

	log.Printf("Preparing insert floor [%d, %d]\n", firstFloorID, lastFloorID)

	res, err := ES.Bulk(BulkBuffer, ES.Bulk.WithIndex(IndexName), ES.Bulk.WithRefresh("wait_for"))
	if err != nil || res.IsError() {
		return fmt.Errorf("error indexing floor [%d, %d]: %s", firstFloorID, lastFloorID, err)
	}
	_ = res.Body.Close()
	log.Printf("index floor [%d, %d] success\n", firstFloorID, lastFloorID)

	BulkBuffer.Reset()
	return nil
}

// BulkDelete used when a hole becomes hidden and delete all of its floors
func BulkDelete(floors Floors) {
	// todo
}

// FloorIndex insert or replace a document, used when a floor is created
// see https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-index_.html
func FloorIndex(floor *Floor) {
	var buffer = bytes.NewBuffer(make([]byte, 16384))

	floorModel := FloorModel{Content: floor.Content}
	err := json.NewEncoder(buffer).Encode(floorModel)
	if err != nil {
		log.Printf("floor encode error: floor_id: %v", floor.ID)
	}

	req := esapi.IndexRequest{
		Index:      IndexName,
		DocumentID: strconv.Itoa(floor.ID),
		Body:       buffer,
		Refresh:    "false",
	}

	res, err := req.Do(context.Background(), ES)
	if err != nil || res.IsError() {
		log.Printf("error index floor: %d\n", floor.ID)
	} else {
		log.Printf("index floor success: %d\n", floor.ID)
	}
}

// FloorDelete used when a floor is deleted
func FloorDelete(floor *Floor) {
	rsp, err := ES.Delete(
		IndexName,
		strconv.Itoa(floor.ID))

	if err != nil || rsp.IsError() {
		log.Printf("error index floor: %d\n", floor.ID)
	} else {
		log.Printf("index floor success: %d\n", floor.ID)
	}
}
