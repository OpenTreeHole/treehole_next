package floor

import (
	"github.com/gofiber/fiber/v2"

	"treehole_next/utils"
)

func RegisterRoutes(app fiber.Router) {
	app.Post("/floors/search", SearchFloors)
	app.Get("/floors/search", SearchFloors)

	app.Get("/holes/:id<int>/floors", ListFloorsInAHole)
	app.Get("/floors", ListFloorsOld)
	app.Get("/floors/:id<int>", GetFloor)
	app.Post("/holes/:id<int>/floors", utils.MiddlewareHasAnsweredQuestions, CreateFloor)
	app.Post("/floors", utils.MiddlewareHasAnsweredQuestions, CreateFloorOld)
	app.Put("/floors/:id<int>", ModifyFloor)
	app.Patch("/floors/:id<int>/_modify", ModifyFloor)
	app.Post("/floors/:id<int>/like/:like<int>", ModifyFloorLike)
	app.Delete("/floors/:id<int>", DeleteFloor)

	app.Get("/users/me/floors", ListReplyFloors)

	app.Get("/floors/:id<int>/history", GetFloorHistory)
	app.Post("/floors/:id<int>/restore/:floor_history_id<int>", RestoreFloor)

	app.Post("/config/search", SearchConfig)
	app.Get("/floors/:id<int>/punishment", GetPunishmentHistory)
	app.Get("/floors/:id<int>/user_silence", GetUserSilence)

	app.Get("/floors/_sensitive", ListSensitiveFloors)
	app.Put("/floors/:id<int>/_sensitive", ModifyFloorSensitive)
	app.Patch("/floors/:id<int>/_sensitive", ModifyFloorSensitive)
}
