package hole

import (
	"github.com/gofiber/fiber/v2"

	"treehole_next/utils"
)

func RegisterRoutes(app fiber.Router) {
	app.Get("/divisions/:id<int>/holes", ListHolesByDivision)
	app.Get("/tags/:name/holes", ListHolesByTag)
	app.Get("/users/me/holes", ListHolesByMe)
	app.Get("/holes/:id<int>", GetHole)
	app.Get("/holes", ListHoles)
	app.Get("/holes/_good", ListGoodHoles)
	app.Get("/holes/_sfw", ListSfwHoles)
	app.Post("/divisions/:id/holes", utils.MiddlewareHasAnsweredQuestions, CreateHole)
	app.Post("/holes", utils.MiddlewareHasAnsweredQuestions, CreateHoleOld)
	app.Patch("/holes/:id<int>/_webvpn", ModifyHole)
	app.Patch("/holes/:id<int>", PatchHole)
	app.Put("/holes/:id<int>", ModifyHole)
	app.Delete("/holes/:id<int>", HideHole)
	app.Delete("/holes/:id<int>/_force", DeleteHole)
}
