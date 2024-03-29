package hole

import (
	"github.com/gofiber/fiber/v2"

	"treehole_next/utils"
)

func RegisterRoutes(app fiber.Router) {
	app.Get("/divisions/:id/holes", ListHolesByDivision)
	app.Get("/tags/:name/holes", ListHolesByTag)
	app.Get("/users/me/holes", ListHolesByMe)
	app.Get("/holes/:id<int>", GetHole)
	app.Get("/holes", ListHolesOld)
	app.Get("/holes/_good", ListGoodHoles)
	app.Post("/divisions/:id/holes", utils.MiddlewareHasAnsweredQuestions, CreateHole)
	app.Post("/holes", utils.MiddlewareHasAnsweredQuestions, CreateHoleOld)
	app.Patch("/holes/:id", PatchHole)
	app.Put("/holes/:id", ModifyHole)
	app.Delete("/holes/:id", HideHole)
	app.Delete("/holes/:id/_force", DeleteHole)
}
