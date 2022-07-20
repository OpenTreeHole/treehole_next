package models

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	BaseModel
	Favorites   []Hole         `json:"favorites" gorm:"many2many:user_favorites"`
	Claims      Map            `json:"claims" gorm:"-:all"`
	Config      Map            `json:"config" gorm:"-:all"`
	BanDivision map[int]bool   `json:"-" gorm:"-:all"`
	Permission  PermissionType `json:"permission" gorm:"-:all"`
}

// PermissionType enum
const (
	P_ADMIN = 1 << iota
	P_OPERATOR
)

type PermissionType uint64

func (user *User) GetUser(c *fiber.Ctx) error {
	id, err := GetUserID(c)
	if err != nil {
		return err
	}
	user.ID = id
	if config.Config.Debug {
		user.Permission = math.MaxUint64
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
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("jwt parse err")
	}
	user.Claims = Map(claims)
	roles, ok := user.Claims["roles"].([]string)
	if !ok {
		return errors.New("jwt parse err")
	}

	for _, v := range roles {
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
