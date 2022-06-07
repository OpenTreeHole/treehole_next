package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"treehole_next/apis"
	"treehole_next/config"
	"treehole_next/middlewares"
	"treehole_next/models"
	"treehole_next/utils"
)

func Init() *fiber.App {
	config.InitConfig()
	models.InitDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.MyErrorHandler,
	})
	middlewares.RegisterMiddlewares(app)
	apis.RegisterRoutes(app)

	return app
}
