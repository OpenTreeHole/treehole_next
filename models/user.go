package models

import (
	"strconv"
	"treehole_next/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/exp/slices"
)

type User struct {
	BaseModel
	Favorites []Hole                 `json:"favorites" gorm:"many2many:user_favorites"`
	Roles     []string               `json:"-" gorm:"-:all"`
	Nickname  string                 `json:"nickname" gorm:"-:all"`
	Config    map[string]interface{} `json:"config" gorm:"-:all"`
	IsAdmin   bool                   `json:"is_admin" gorm:"-:all"`
}

func (user *User) GetUser(c *fiber.Ctx) error {
	id, err := GetUserID(c)
	if err != nil {
		return err
	}
	if config.Config.Debug {
		user.IsAdmin = true
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	user.Roles = claims["role"].([]string)
	user.Nickname = claims["nickname"].(string)
	if slices.Contains(user.Roles, "admin") {
		user.IsAdmin = true
	}

	user.ID = id
	return nil
}

func GetUserID(c *fiber.Ctx) (int, error) {
	if config.Config.Debug {
		return 1, nil
	}

	id, err := strconv.Atoi(c.Get("X-Consumer-Username"))
	if err != nil {
		return 0, err
	}

	return id, nil
}
