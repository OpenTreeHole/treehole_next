package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm/clause"
	. "treehole_next/models"
)

func RegisterRoutes(app fiber.Router) {
	app.Get("/users/me", GetCurrentUser)
	app.Get("/users/:id", GetUserByID)
	app.Put("/users/:id", ModifyUser)
}

// GetCurrentUser
//
// @Summary get current user
// @Tags user
// @Deprecated
// @Produce json
// @Router /users/me [get]
// @Success 200 {object} User
func GetCurrentUser(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return err
	}
	return c.JSON(&user)
}

// GetUserByID
//
// @Summary get user by id, owner or admin
// @Tags user
// @Produce json
// @Router /users/{user_id} [get]
// @Success 200 {object} User
func GetUserByID(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	user, err := GetUser(c)
	if err != nil {
		return err
	}

	if !user.IsAdmin || user.ID == userID {
		return common.Forbidden()
	}

	var getUser User
	err = getUser.LoadUserByID(userID)
	if err != nil {
		return err
	}

	return c.JSON(&getUser)
}

// ModifyUser
//
// @Summary modify user profiles
// @Tags User
// @Produce json
// @Router /users/{user_id} [put]
// @Success 200 {object} User
func ModifyUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	user, err := GetUser(c)
	if err != nil {
		return err
	}

	if !user.IsAdmin && user.ID != userID {
		return common.Forbidden()
	}

	body, err := common.ValidateBody[ModifyModel](c)
	if err != nil {
		return err
	}

	if body.Config != nil {
		user.Config = *body.Config
	}

	err = DB.Model(&user).Omit(clause.Associations).Select("Config").Updates(&user).Error
	if err != nil {
		return err
	}

	return c.JSON(&user)
}
