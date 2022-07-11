package utils

import (
	"encoding/json"
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
