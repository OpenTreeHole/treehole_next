package favourite

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/user/favorites", ListFavorites)
	app.Post("/user/favorites", AddFavorite)
	app.Put("/user/favorites", ModifyFavorite)
	app.Delete("/user/favorites", DeleteFavorite)
	app.Get("/user/favorite_group", ListFavoriteGroups)
	app.Post("/user/favorite_group", AddFavoriteGroup)
	app.Put("/user/favorite_group", ModifyFavoriteGroup)
	app.Delete("/user/favorite_group", DeleteFavoriteGroup)
	app.Put("/user/favorite_group/move", MoveFavorite)
}
