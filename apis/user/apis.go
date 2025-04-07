package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm/clause"

	. "treehole_next/models"
)

func RegisterRoutes(app fiber.Router) {
	app.Get("/users/me", GetCurrentUser)
	app.Get("/users/:id<int>", GetUserByID)
	app.Put("/users/:id<int>", ModifyUser)
	app.Patch("/users/:id<int>/_webvpn", ModifyUser)
	app.Put("/users/me", ModifyCurrentUser)
	app.Patch("/users/me/_webvpn", ModifyCurrentUser)
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
	user, err := GetCurrLoginUser(c)
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
// @Param user_id path int true "user id"
// @Success 200 {object} User
func GetUserByID(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	user, err := GetCurrLoginUser(c)
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
// @Router /users/{user_id}/_webvpn [patch]
// @Param user_id path int true "user id"
// @Param json body ModifyModel true "modify user"
// @Success 200 {object} User
func ModifyUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	user, err := GetCurrLoginUser(c)
	if err != nil {
		return err
	}

	if !user.IsAdmin && user.ID != userID {
		return common.Forbidden()
	}

	var body ModifyModel
	err = common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// cannot get field "has_answered_questions" when admin changes other user's config
	if user.ID != userID {
		user = &User{
			ID: userID,
		}
		err = DB.Take(user).Error
		if err != nil {
			return err
		}
	}

	err = modifyUser(c, user, body)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// ModifyCurrentUser
//
// @Summary modify current user profiles
// @Tags User
// @Produce json
// @Router /users/me [put]
// @Router /users/me/_webvpn [patch]
// @Param user_id path int true "user id"
// @Param json body ModifyModel true "modify user"
// @Success 200 {object} User
func ModifyCurrentUser(c *fiber.Ctx) error {
	user, err := GetCurrLoginUser(c)
	if err != nil {
		return err
	}

	var body ModifyModel
	err = common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	err = modifyUser(c, user, body)
	if err != nil {
		return err
	}

	return c.JSON(&user)
}

func modifyUser(_ *fiber.Ctx, user *User, body ModifyModel) error {
	var newUser User
	err := DB.Select("config").First(&newUser, user.ID).Error
	if err != nil {
		return err
	}

	if body.Config != nil {
		if body.Config.Notify != nil {
			newUser.Config.Notify = body.Config.Notify
		}
		if body.Config.ShowFolded != nil {
			newUser.Config.ShowFolded = *body.Config.ShowFolded
		}
	}

	err = DB.Model(&user).Omit(clause.Associations).Select("Config").UpdateColumns(&newUser).Error
	if err != nil {
		return err
	}

	user.Config = newUser.Config
	return nil
}
