package models

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"treehole_next/config"
)

type User struct {
	BaseModel
	Favorites []Hole    `json:"favorites" gorm:"many2many:user_favorites"`
	Nickname  string    `json:"nickname" gorm:"-:all"`
	Config    StringMap `json:"config" gorm:"-:all"`
	IsAdmin   bool      `json:"is_admin" gorm:"-:all"`
}

func (user *User) GetUser(c *fiber.Ctx) error {
	if config.Config.Debug {
		user.ID = 1
		return nil
	}

	id, err := strconv.Atoi(c.Get("X-Consumer-Username"))
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}
