package division

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Post("/divisions", AddDivision)
	app.Get("/divisions", ListDivisions)
	app.Get("/divisions/:id", GetDivision)
	app.Put("/divisions/:id", ModifyDivision)
	app.Patch("/divisions/:id/_webvpn", ModifyDivision)
	app.Delete("/divisions/:id", DeleteDivision)
}
