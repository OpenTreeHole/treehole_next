package bootstrap

import (
	"github.com/goccy/go-json"
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
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	middlewares.RegisterMiddlewares(app)
	apis.RegisterRoutes(app)

	startTasks()

	return app
}
