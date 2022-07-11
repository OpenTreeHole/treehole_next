package utils

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// BindJSON is a safe method to bind request body to struct
func BindJSON(c *fiber.Ctx, obj interface{}) error {
	body := c.Body()
	if len(body) == 0 {
		body, _ = json.Marshal(fiber.Map{})
	}
	return json.Unmarshal(body, obj)
}

type CanPreprocess interface {
	Preprocess() error
}

func Serialize(c *fiber.Ctx, obj CanPreprocess) error {
	err := obj.Preprocess()
	if err != nil {
		return err
	}
	return c.JSON(obj)
}

func ReText2IntArray(IDs [][]string) ([]int, error) {
	ansIDMapping := make(map[int]bool)
	for _, v := range IDs {
		id, err := strconv.Atoi(v[1])
		if err != nil {
			return nil, err
		}
		ansIDMapping[id] = true
	}
	keys := []int{}
	for key := range ansIDMapping {
		keys = append(keys, key)
	}
	return keys, nil
}
