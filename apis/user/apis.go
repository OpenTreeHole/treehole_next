package user

import (
	"github.com/gofiber/fiber/v2"
	. "treehole_next/models"
	"treehole_next/utils"
)

func RegisterRoutes(app fiber.Router) {
	app.Get("/users/me", GetCurrentUser)
	app.Get("/users/:id", GetUserByID)
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
// @Deprecated
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
		return utils.Forbidden()
	}

	var getUser User
	err = getUser.LoadUserByID(userID)
	if err != nil {
		return err
	}

	return c.JSON(&getUser)
}
