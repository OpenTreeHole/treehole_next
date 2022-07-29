package search

import (
	"bytes"
	"encoding/json"
	. "treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

// SearchFloors
// @Summary SearchFloors In Elastic Search
// @Tags Search
// @Produce application/json
// @Router /floors/search [post]
// @Param json body any true "json"
// @Success 200 {array} models.Floor
func SearchFloors(c *fiber.Ctx) error {
	// forwarding
	var reqBody bytes.Buffer
	reqBody.Write(c.Body())
	res, err := ES.Search(
		ES.Search.WithIndex("floor"),
		ES.Search.WithBody(&reqBody),
	)
	if err != nil {
		return err
	}
	if res.IsError() {
		e := Map{}
		err := json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return err
		} else {
			return c.Status(502).JSON(&e)
		}
	}

	var resBody Map
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return err
	}
	floorIDs := make([]int, 0, 20)
	for _, hit := range resBody["hits"].(Map)["hits"].(Map) {
		floorIDs = append(floorIDs, hit.(Map)["_id"].(int))
	}

	// get floors
	var floors Floors
	result := DB.Preload("Mention").Find(floors, floorIDs)
	if result.Error != nil {
		return result.Error
	}

	// order
	var orderedFloors Floors
	for _, order := range floorIDs {
		index := func(target int) int {
			left := 0
			right := len(floors)
			for left < right {
				mid := left + (right-left)>>1
				if floors[mid].ID < target {
					left = mid + 1
				} else if floors[mid].ID > target {
					right = mid
				} else {
					return mid
				}
			}
			return -1
		}(order)
		if index >= 0 {
			orderedFloors = append(orderedFloors, floors[index])
		}
	}
	return Serialize(c, orderedFloors)
}
