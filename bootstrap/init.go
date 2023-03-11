package bootstrap

import (
	"context"
	"time"
	"treehole_next/apis"
	"treehole_next/apis/hole"
	"treehole_next/apis/message"
	"treehole_next/config"
	"treehole_next/models"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func Init() (*fiber.App, context.CancelFunc) {
	config.InitConfig()
	utils.InitCache()
	models.Init()
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
		app.Use(MyLogger)
	}
	if config.Config.Mode == "dev" {
		app.Use(pprof.New())
	}
	app.Use(GetUser)
}

func GetUser(c *fiber.Ctx) error {
	user, err := models.GetUser(c)
	if err != nil {
		return err
	}
	c.Locals("user", user)

	return c.Next()
}

func MyLogger(c *fiber.Ctx) error {
	startTime := time.Now()
	err := c.Next()
	latency := time.Since(startTime)
	user := c.Locals("user").(*models.User)
	utils.Logger.Info("LOG : ",
		zap.Int("StatusCode", c.Response().StatusCode()),
		zap.String("Method", string(c.Context().Method())),
		zap.String("OriginUrl", c.OriginalURL()),
		zap.String("RemoteIP", string(c.Context().RemoteIP())),
		zap.Int("Latency", int(latency)),
		zap.String("Error", err.Error()),
		zap.Int("User", user.ID),
	)
	return err
}

func startTasks() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go hole.UpdateHoleViews(ctx)
	go message.PurgeMessage()
	go models.UpdateAdminList(ctx)
	return cancel
}
