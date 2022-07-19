package middlewares

import (
	"treehole_next/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RegisterMiddlewares(app *fiber.App) {
	app.Use(recover.New())
	if config.Config.Mode != "perf" {
		app.Use(logger.New())
	}
}
