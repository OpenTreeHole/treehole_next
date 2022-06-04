package apis

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"treehole_next/apis/division"
	_ "treehole_next/docs"
)

func registerRoutes(app *fiber.App) {
	app.Get("/", index)
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)
}

func RegisterRoutes(app *fiber.App) {
	registerRoutes(app)
	division.RegisterRoutes(app)
}
