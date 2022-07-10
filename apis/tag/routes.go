package tag

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	app.Get("/tags", ListTags)
	app.Get("/tags/:id", GetTag)
	app.Post("/tags", CreateTag)
	app.Put("/tags/:id", ModifyTag)
	app.Delete("/tags/:id", DeleteTag)
}
