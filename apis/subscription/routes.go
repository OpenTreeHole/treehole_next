package subscription

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/user/subscriptions", ListSubscriptions)
	app.Post("/user/subscriptions", AddSubscription)
	app.Delete("/user/subscription", DeleteSubscription)
}
