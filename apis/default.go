package apis

import "github.com/gofiber/fiber/v2"

// Index
// @Produce application/json
// @Success 200 {object} utils.MessageModel
// @Router / [get]
func Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "hello world"})
}
