package apis

import (
	"github.com/opentreehole/go-common"

	"treehole_next/apis/division"
	"treehole_next/apis/favourite"
	"treehole_next/apis/floor"
	"treehole_next/apis/hole"
	"treehole_next/apis/message"
	"treehole_next/apis/penalty"
	"treehole_next/apis/report"
	"treehole_next/apis/subscription"
	"treehole_next/apis/tag"
	"treehole_next/apis/user"
	"treehole_next/config"
	_ "treehole_next/docs"
	"treehole_next/models"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func registerRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)
}

func RegisterRoutes(app *fiber.App) {
	registerRoutes(app)

	group := app.Group("/api")
	group.Get("/", Index)
	group.Use(MiddlewareGetUser)
	division.RegisterRoutes(group)
	tag.RegisterRoutes(group)
	hole.RegisterRoutes(group)
	floor.RegisterRoutes(group)
	report.RegisterRoutes(group)
	favourite.RegisterRoutes(group)
	subscription.RegisterRoutes(group)
	penalty.RegisterRoutes(group)
	user.RegisterRoutes(group)
	message.RegisterRoutes(group)
}

func MiddlewareGetUser(c *fiber.Ctx) error {
	userObject, err := models.GetUser(c)
	if err != nil {
		return err
	}
	c.Locals("user", userObject)
	if config.Config.AdminOnly {
		if !userObject.IsAdmin {
			return common.Forbidden()
		}
	}
	return c.Next()
}
