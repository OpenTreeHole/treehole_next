package floor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	. "treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"
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
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// SearchFloors
// @Summary SearchFloors In ElasticSearch
// @Tags Search
// @Produce application/json
// @Router /floors/search [post]
// @Param json body any true "json"
// @Success 200 {array} models.Floor
func SearchFloors(c *fiber.Ctx) error {
	// forwarding
	var body bytes.Buffer
	body.Write(c.Body())
	return search(c, body)
}

func SearchFloorsOld(c *fiber.Ctx, query ListOldModel) error {
	var body bytes.Buffer
	body.WriteString(fmt.Sprintf(`
	{
		"query": {
			"match": {
				"content": {
					"query": "%s",
					"operator": "or"
				}
			}
		},
		"sort": {
			"_score": {
				"order": "desc"
			}
		},
		"size": %d,
		"from": %d
	}`, query.Search, query.Size, query.Offset))
	return search(c, body)
}

func search(c *fiber.Ctx, body bytes.Buffer) error {
	res, err := ES.Search(
		ES.Search.WithIndex("floor"),
		ES.Search.WithBody(&body),
	)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	if res.IsError() {
		e := Map{}
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return err
		} else {
			return c.Status(502).JSON(&e)
		}
	}

	var response SearchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	floorIDs := make([]int, len(response.Hits.Hits))
	for i, hit := range response.Hits.Hits {
		floorIDs[i] = hit.Source.ID
	}
	Logger.Debug("search response", zap.Ints("floorIDs", floorIDs))

	// get floors
	var floors Floors
	result := DB.Preload("Mention").Find(&floors, floorIDs)
	if result.Error != nil {
		return result.Error
	}

	// order
	floors = OrderInGivenOrder(floors, floorIDs)
	return Serialize(c, floors)
}
