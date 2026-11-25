package apis

import (
	"github.com/opentreehole/go-common"
	"strconv"
	"strings"

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
	group.Use(MiddlewareContentAccess)
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

// MiddlewareContentAccess guards direct URLs, including nested floor and
// summary routes. List queries are filtered separately before pagination.
func MiddlewareContentAccess(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Next()
	}
	enabled, err := models.ContentFilterEnabled(c)
	if err != nil {
		return err
	}
	if !enabled {
		return c.Next()
	}
	user, err := models.GetCurrLoginUser(c)
	if err != nil {
		return err
	}
	parts := strings.Split(strings.Trim(c.Path(), "/"), "/")
	if len(parts) < 3 || parts[0] != "api" {
		return c.Next()
	}
	var count int64
	switch parts[1] {
	case "divisions":
		id, parseErr := strconv.Atoi(parts[2])
		if parseErr != nil {
			return c.Next()
		}
		err = models.DB.Model(&models.UserHiddenDivision{}).Where("user_id = ? AND division_id = ?", user.ID, id).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return common.NotFound()
		}
	case "holes", "floors":
		id, parseErr := strconv.Atoi(parts[2])
		if parseErr != nil {
			return c.Next()
		}
		query := models.FilterHoles(models.DB.Model(&models.Hole{}), user.ID)
		if parts[1] == "holes" {
			query = query.Where("hole.id = ?", id)
		} else {
			query = query.Joins("JOIN floor ON floor.hole_id = hole.id").Where("floor.id = ?", id)
		}
		if err = query.Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return common.NotFound()
		}
	case "tags":
		query := models.DB.Model(&models.UserHiddenTag{}).Where("user_id = ?", user.ID)
		if id, parseErr := strconv.Atoi(parts[2]); parseErr == nil {
			query = query.Where("LOWER(name) IN (?)", models.DB.Model(&models.Tag{}).Select("LOWER(name)").Where("id = ?", id))
		} else {
			query = query.Where("LOWER(name) = LOWER(?)", parts[2])
		}
		if err = query.Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return common.NotFound()
		}
	}
	return c.Next()
}

func MiddlewareGetUser(c *fiber.Ctx) error {
	userObject, err := models.GetCurrLoginUser(c)
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
