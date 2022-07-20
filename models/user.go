package models

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	BaseModel
	Favorites   []Hole         `json:"favorites" gorm:"many2many:user_favorites"`
	PayloadMap  Map            `json:"-" gorm:"-:all"`
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
	tokenString := rawToken[7:]                       // extract "Bearer "
	payload := strings.SplitN(tokenString, ".", 2)[1] // the middle one is payload
	payloadBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return utils.BadRequest("invalid jwt token. Fail to decode payload.")
	}
	payloadMap := Map{}
	err = json.Unmarshal(payloadBytes, &payloadMap)
	if err != nil {
		return utils.BadRequest("invalid jwt token. Fail to unmarchal payload.")
	}
	user.PayloadMap = payloadMap
	roles, ok := user.PayloadMap["roles"].([]string)
	if !ok {
		return utils.BadRequest("invalid jwt token. Fail to parse roles")
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
