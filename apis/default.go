package apis

import (
	"github.com/gofiber/fiber/v2"
)

// Index
// @Produce application/json
// @Success 200 {object} models.MessageModel
// @Router / [get]
func Index(c *fiber.Ctx) error {
	return c.SendFile("data/meta.json")
}
