package message

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Post("/messages", SendMail)
	app.Get("/messages", ListMessages)
	app.Post("/messages/clear", ClearMessages)
	app.Put("/messages", ClearMessagesDeprecated)
	app.Patch("/messages/_clear", ClearMessagesDeprecated)
	app.Delete("/messages/:id<int>", DeleteMessage)
}
