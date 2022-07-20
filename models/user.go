package models

import (
	"errors"
	"strconv"
	"strings"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	BaseModel
	Favorites   []Hole                 `json:"favorites" gorm:"many2many:user_favorites"`
	Roles       []string               `json:"-" gorm:"-:all"`
	BanDivision map[int]bool           `json:"-" gorm:"-:all"`
	Nickname    string                 `json:"nickname" gorm:"-:all"`
	Config      map[string]interface{} `json:"config" gorm:"-:all"`
	Permission  PermissionType         `json:"permission" gorm:"-:all"`
}

// PermissionType enum
const (
	P_ADMIN = 1 << iota
	P_OPERATOR
)

type PermissionType int

func (user *User) GetUser(c *fiber.Ctx) error {
	id, err := GetUserID(c)
	if err != nil {
		return err
	}
	user.ID = id
	if config.Config.Debug {
		user.Permission = P_ADMIN + P_OPERATOR
		return nil
	}

	// extract and parse token
	rawToken := c.Get("Authorization")
	tokenString := rawToken[7:] // extract "Bearer "
	userToken, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	// get userinfo
	claims := userToken.Claims.(jwt.MapClaims)
	roles, ok := claims["roles"].([]string)
	if !ok {
		return errors.New("jwt parse err")
	}
	user.Roles = roles
	nickname, ok := claims["nickname"].(string)
	if !ok {
		return errors.New("jwt parse err")
	}
	user.Nickname = nickname
	for _, v := range user.Roles {
		if v == "admin" {
			user.Permission |= P_ADMIN
		} else if v == "operator" {
			user.Permission |= P_OPERATOR
		} else if strings.HasPrefix(v, "ban_treehole") {
			banDivisionID, err := strconv.Atoi(v[13:])
			if err != nil {
				return err
			}
			user.BanDivision[banDivisionID] = true
		}
	}
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

// get userInfo and check user permission
//
// Example:
//
//  GetAndCheckPermission(c, P_ADMIN | P_OPERATOR)
//
func (user *User) GetAndCheckPermission(c *fiber.Ctx, t PermissionType) error {
	err := user.GetUser(c)
	if err != nil {
		return err
	}
	if !user.CheckPermission(t) {
		return utils.Forbidden()
	}
	return nil
}

// check user permission
//
// Example:
//
//  CheckPermission(P_ADMIN)
//  CheckPermission(P_ADMIN | P_OPERATOR)
//
func (user *User) CheckPermission(t PermissionType) bool {
	return user.Permission&t != 0
}
