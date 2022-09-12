package apis

import (
	"treehole_next/data"

	"github.com/gofiber/fiber/v2"
)

// Index
// @Produce application/json
// @Success 200 {object} models.MessageModel
// @Router / [get]
func Index(c *fiber.Ctx) error {
	return c.Send(data.MetaFile)
}
