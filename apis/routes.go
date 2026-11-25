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
	group.Use("/divisions/:id<int>", visibleResource("division"))
	group.Use("/holes/:id<int>", visibleResource("hole"))
	group.Use("/floors/:id<int>", visibleResource("floor"))
	division.RegisterRoutes(group)
	tag.RegisterRoutes(group, visibleResource("tag"))
	hole.RegisterRoutes(group)
	floor.RegisterRoutes(group)
	report.RegisterRoutes(group)
	favourite.RegisterRoutes(group)
	subscription.RegisterRoutes(group)
	penalty.RegisterRoutes(group)
	user.RegisterRoutes(group)
	message.RegisterRoutes(group)
}

// Use router parameters so aliases and every HTTP method share the same check.
func visibleResource(kind string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, err := models.FilterUserID(c)
		if err != nil {
			return err
		}
		if uid == 0 {
			return c.Next()
		}
		id, err := c.ParamsInt("id")
		if err != nil {
			return common.NotFound()
		}
		var count int64
		switch kind {
		case "division":
			// Zero is the legacy all-divisions list.
			if id != 0 {
				err = models.RequireVisibleDivision(models.DB, uid, id)
			}
		case "hole":
			err = models.RequireVisibleHoles(models.DB, uid, []int{id})
		case "floor":
			err = models.VisibleFloorIDs(models.DB, uid).Where("floor.id = ?", id).Count(&count).Error
			if err == nil && count == 0 {
				return common.NotFound()
			}
		case "tag":
			err = models.SelectableTags(models.DB, uid).Where("tag.id = ?", id).Count(&count).Error
			if err == nil && count == 0 {
				return common.NotFound()
			}
		}
		if err != nil {
			return err
		}
		return c.Next()
	}
}

func MiddlewareGetUser(c *fiber.Ctx) error {
	userObject, err := models.GetCurrLoginUser(c)
	if err != nil {
		return err
	}
	c.Locals("user", userObject)
	// Resolve the request policy before a handler opens a write transaction.
	if _, err := models.ContentFilterEnabled(c); err != nil {
		return err
	}
	if config.Config.AdminOnly {
		if !userObject.IsAdmin {
			return common.Forbidden()
		}
	}
	return c.Next()
}
