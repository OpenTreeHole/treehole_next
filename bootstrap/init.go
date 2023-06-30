package bootstrap

import (
	"context"
	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"
	"time"
	"treehole_next/apis"
	"treehole_next/apis/hole"
	"treehole_next/apis/message"
	"treehole_next/config"
	"treehole_next/models"
	"treehole_next/utils"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Init() (*fiber.App, context.CancelFunc) {
	config.InitConfig()
	utils.InitCache()
	models.Init()
	models.InitDB()
	models.InitAdminList()

	app := fiber.New(fiber.Config{
		ErrorHandler: common.CommonErrorHandler,
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
	app.Use(pprof.New())
	app.Use(GetUserID)
}

func GetUserID(c *fiber.Ctx) error {
	userID, err := models.GetUserID(c)
	if err == nil {
		c.Locals("user_id", userID)
	}

	return c.Next()
}

func MyLogger(c *fiber.Ctx) error {
	startTime := time.Now()
	chainErr := c.Next()

	if chainErr != nil {
		if err := c.App().ErrorHandler(c, chainErr); err != nil {
			_ = c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	latency := time.Since(startTime).Milliseconds()
	userID, ok := c.Locals("user_id").(int)

	output := log.Info().
		Int("status_code", c.Response().StatusCode()).
		Str("method", c.Method()).
		Str("origin_url", c.OriginalURL()).
		Str("remote_ip", c.Get("X-Real-IP")).
		Int64("latency", latency)
	if ok {
		output = output.Int("user_id", userID)
	}
	if chainErr != nil {
		output = output.Err(chainErr)
	}
	output.Msg("http log")
	return nil
}

func startTasks() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go hole.UpdateHoleViews(ctx)
	go message.PurgeMessage()
	go models.UpdateAdminList(ctx)
	return cancel
}
