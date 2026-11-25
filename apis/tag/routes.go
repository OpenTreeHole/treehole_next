package tag

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router, visibleTag fiber.Handler) {
	app.Get("/tags", ListTags)
	app.Get("/tags/:id<int>", visibleTag, GetTag)
	app.Post("/tags", CreateTag)
	app.Put("/tags/:id<int>", visibleTag, ModifyTag)
	app.Patch("/tags/:id<int>/_webvpn", visibleTag, ModifyTag)
	app.Delete("/tags/:id<int>", visibleTag, DeleteTag)
}
