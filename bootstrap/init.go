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
	output := []zap.Field{
		zap.Int("StatusCode", c.Response().StatusCode()),
		zap.String("Method", c.Method()),
		zap.String("OriginUrl", c.OriginalURL()),
		zap.String("RemoteIP", c.Get("X-Real-IP")),
		zap.Int64("Latency", latency),
	}
	if ok {
		output = append(output, zap.Int("UserID", userID))
	}
	if chainErr != nil {
		output = append(output, zap.Error(chainErr))
	}
	utils.Logger.Info("LOG : ", output...)
	return nil
}

func startTasks() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go hole.UpdateHoleViews(ctx)
	go message.PurgeMessage()
	go models.UpdateAdminList(ctx)
	return cancel
}
