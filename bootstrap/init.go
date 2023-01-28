package bootstrap

import (
	"context"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"treehole_next/apis"
	"treehole_next/apis/hole"
	"treehole_next/config"
	"treehole_next/models"
	"treehole_next/utils"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func Init() (*fiber.App, chan struct{}) {
	models.InitDB()
	utils.Logger, _ = utils.InitLog()
	models.InitAdminList()

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.MyErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	registerMiddlewares(app)
	apis.RegisterRoutes(app)

	return app, startTasks()
}

func registerMiddlewares(app *fiber.App) {
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	if config.Config.Mode != "bench" {
		app.Use(logger.New())
	}
	if config.Config.Mode == "dev" {
		app.Use(pprof.New())
	}
}

func startTasks() chan struct{} {
	done := make(chan struct{}, 1)
	go hole.UpdateHoleViews(done)
	go models.UpdateTagTemperature(context.Background())
	return done
}
