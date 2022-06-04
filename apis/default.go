package apis

import "github.com/gofiber/fiber/v2"

// index
// @Produce application/json
// @Success 200 {object} utils.MessageModel
// @Router / [get]
func index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "hello world"})
}
