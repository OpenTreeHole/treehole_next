package elastic

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/goccy/go-json"
	"io"
	"log"
	"strconv"
	"treehole_next/utils"

	. "treehole_next/models"
)

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
