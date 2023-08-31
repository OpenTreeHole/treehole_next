package floor

import (
	"github.com/gofiber/fiber/v2"

	"treehole_next/utils"
)

func RegisterRoutes(app fiber.Router) {
	app.Post("/floors/search", SearchFloors)
	app.Get("/floors/search", SearchFloors)

	app.Get("/holes/:id/floors", ListFloorsInAHole)
	app.Get("/floors", ListFloorsOld)
	app.Get("/floors/:id", GetFloor)
	app.Post("/holes/:id/floors", utils.MiddlewareHasAnsweredQuestions, CreateFloor)
	app.Post("/floors", utils.MiddlewareHasAnsweredQuestions, CreateFloorOld)
	app.Put("/floors/:id", ModifyFloor)
	app.Post("/floors/:id/like/:like", ModifyFloorLike)
	app.Delete("/floors/:id", DeleteFloor)

	app.Get("/users/me/floors", ListReplyFloors)

	app.Get("/floors/:id/history", GetFloorHistory)
	app.Post("/floors/:id/restore/:floor_history_id", RestoreFloor)

	app.Post("/config/search", SearchConfig)
	app.Get("/floors/:id/punishment", GetPunishmentHistory)
}
