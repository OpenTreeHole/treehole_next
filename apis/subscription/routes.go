package subscription

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/users/subscriptions", ListSubscriptions)
	app.Post("/users/subscriptions", AddSubscription)
	app.Delete("/users/subscription", DeleteSubscription)
}
