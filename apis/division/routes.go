package division

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	app.Post("/divisions", AddDivision)
	app.Get("/divisions", ListDivisions)
	app.Get("/divisions/:id", GetDivision)
	app.Put("/divisions/:id", ModifyDivision)
	app.Delete("/divisions/:id", DeleteDivision)
}
