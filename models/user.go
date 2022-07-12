package models

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"treehole_next/config"
)

type User struct {
	BaseModel
	Favorites []Hole                 `json:"favorites" gorm:"many2many:user_favorites"`
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
	// TODO: jwt
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
