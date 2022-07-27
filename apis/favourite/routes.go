package favourite

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/user/favorites", ListFavorites)
	app.Post("/user/favorites", AddFavorite)
	app.Put("/user/favorites", ModifyFavorite)
	app.Delete("/user/favorites", DeleteFavorite)
}
