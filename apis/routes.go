package apis

import (
	_ "treehole_next/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
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

	// holes
	app.Get("/divisions/:id/holes", ListHolesByDivision)
	app.Get("/tags/:name/holes", ListHolesByTag)
	app.Get("/holes/:id", GetHole)
	app.Get("/holes", ListHolesOld)
	app.Post("/divisions/:id/holes", CreateHole)
	app.Post("/holes", CreateHoleOld)
	app.Put("/holes/:id", ModifyHole)
	app.Delete("/holes/:id", DeleteHole)

	// floors
	app.Get("/holes/:id/floors", ListFloorsInAHole)
	app.Get("/floors", ListFloorsOld)
	app.Get("/floors/:id", GetFloor)
	app.Post("/holes/:id/floors", CreateFloor)
	app.Post("/floors", CreateFloorOld)
	app.Put("/floors/:id", ModifyFloor)
	app.Delete("/floors/:id", DeleteFloor)
}
