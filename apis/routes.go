package apis

import (
	"treehole_next/apis/division"
	"treehole_next/apis/favourite"
	"treehole_next/apis/floor"
	"treehole_next/apis/hole"
	"treehole_next/apis/penalty"
	"treehole_next/apis/report"
	"treehole_next/apis/tag"
	_ "treehole_next/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func registerRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)
}

func RegisterRoutes(app *fiber.App) {
	registerRoutes(app)

	group := app.Group("/api")
	group.Get("/", Index)
	division.RegisterRoutes(group)
	tag.RegisterRoutes(group)
	hole.RegisterRoutes(group)
	floor.RegisterRoutes(group)
	report.RegisterRoutes(group)
	favourite.RegisterRoutes(group)
	penalty.RegisterRoutes(group)
}
