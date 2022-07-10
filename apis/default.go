package apis

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "treehole_next/docs"
)

// Index
// @Produce application/json
// @Success 200 {object} models.MessageModel
// @Router / [get]
func Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "hello world"})
}

func registerRoutes(app *fiber.App) {
	app.Get("/", Index)
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)
}
