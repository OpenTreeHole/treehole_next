package bootstrap

import (
	"treehole_next/apis"
	"treehole_next/config"
	"treehole_next/middlewares"
	"treehole_next/models"
	"treehole_next/utils"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func Init() (*fiber.App, chan struct{}) {
	config.InitConfig()
	models.InitDB()
	utils.Logger, _ = utils.InitLog()
	models.InitAdminList()

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.MyErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	middlewares.RegisterMiddlewares(app)
	apis.RegisterRoutes(app)

	return app, startTasks()
}
