package apis

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "treehole_next/docs"
)

func RegisterRoutes(app *fiber.App) {
	// base
	app.Get("/", index)
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)
	
	// divisions
	app.Post("/divisions", AddDivision)
	app.Get("/divisions", ListDivisions)
	app.Get("/divisions/:id", GetDivision)
}
