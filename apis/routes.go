package apis

import (
	"github.com/gofiber/fiber/v2"
	"treehole_next/apis/division"
	"treehole_next/apis/floor"
	"treehole_next/apis/hole"
	"treehole_next/apis/tag"
)

func RegisterRoutes(app *fiber.App) {
	registerRoutes(app)
	division.RegisterRoutes(app)
	tag.RegisterRoutes(app)
	hole.RegisterRoutes(app)
	floor.RegisterRoutes(app)
}
