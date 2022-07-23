package bootstrap

import (
	"treehole_next/apis"
	"treehole_next/config"
	"treehole_next/middlewares"
	"treehole_next/models"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

func Init() *fiber.App {
	config.InitConfig()
	models.InitDB()
	utils.Logger, _ = utils.InitLog()

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.MyErrorHandler,
	})
	middlewares.RegisterMiddlewares(app)
	apis.RegisterRoutes(app)

	startTasks()

	return app
}
