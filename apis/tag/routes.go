package tag

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/tags", ListTags)
	app.Get("/tags/:id<int>", GetTag)
	app.Post("/tags", CreateTag)
	app.Put("/tags/:id<int>", ModifyTag)
	app.Patch("/tags/:id<int>/_modify", ModifyTag)
	app.Delete("/tags/:id<int>", DeleteTag)
}
