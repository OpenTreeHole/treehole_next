package favourite

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/user/favorites", ListFavorites)
	app.Post("/user/favorites", AddFavorite)
	app.Put("/user/favorites", ModifyFavorite)
	app.Patch("/user/favorites/_webvpn", ModifyFavorite)
	app.Delete("/user/favorites", DeleteFavorite)
	app.Get("/user/favorite_groups", ListFavoriteGroups)
	app.Post("/user/favorite_groups", AddFavoriteGroup)
	app.Put("/user/favorite_groups", ModifyFavoriteGroup)
	app.Patch("/user/favorite_groups/_webvpn", ModifyFavoriteGroup)
	app.Delete("/user/favorite_groups", DeleteFavoriteGroup)
	app.Put("/user/favorites/move", MoveFavorite)
}
