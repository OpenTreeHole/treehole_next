package hole

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/divisions/:id/holes", ListHolesByDivision)
	app.Get("/tags/:name/holes", ListHolesByTag)
	app.Get("/users/me/holes", ListHolesByMe)
	app.Get("/holes/:id", GetHole)
	app.Get("/holes", ListHolesOld)
	app.Post("/divisions/:id/holes", CreateHole)
	app.Post("/holes", CreateHoleOld)
	app.Patch("/holes/:id", PatchHole)
	app.Put("/holes/:id", ModifyHole)
	app.Delete("/holes/:id", HideHole)
	app.Delete("/holes/:id/_force", DeleteHole)
}
