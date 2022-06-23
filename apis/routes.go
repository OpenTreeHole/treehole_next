package apis

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "treehole_next/docs"
)

func RegisterRoutes(app *fiber.App) {
	// base
	app.Get("/", Index)
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	// divisions
	app.Post("/divisions", AddDivision)
	app.Get("/divisions", ListDivisions)
	app.Get("/divisions/:id", GetDivision)
	app.Put("/divisions/:id", ModifyDivision)
	app.Delete("/divisions/:id", DeleteDivision)

	// tags
	app.Get("/tags", ListTags)
	app.Get("/tags/:id", GetTag)
	app.Post("/tags", CreateTag)
	app.Put("/tags/:id", ModifyTag)
	app.Delete("/tags/:id", DeleteTag)
}
