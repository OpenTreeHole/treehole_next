package models

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	BaseModel
	Favorites   []Hole                 `json:"favorites" gorm:"many2many:user_favorites"`
	Roles       []string               `json:"roles" gorm:"-:all"`
	BanDivision map[int]bool           `json:"-" gorm:"-:all"`
	Nickname    string                 `json:"nickname" gorm:"-:all"`
	Config      map[string]interface{} `json:"config" gorm:"-:all"`
	Permission  PermissionType         `json:"permission" gorm:"-:all"`
}

// PermissionType enum
//goland:noinspection GoSnakeCaseUsage
const (
	P_ADMIN = 1 << iota
	P_OPERATOR
)

type PermissionType int

// parseJWT extracts and parse token
func (user *User) parseJWT(token string) bool {
	if len(token) < 7 {
		return false
	}

	payloads := strings.SplitN(token[7:], ".", 3) // extract "Bearer "
	if len(payloads) < 3 {
		return false
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(payloads[1]) // the middle one is payload
	if err != nil {
		return false
	}

	err = json.Unmarshal(payloadBytes, user)
	if err != nil {
		return false
	}

	return true
}

func (user *User) GetUser(c *fiber.Ctx) error {
	if config.Config.Debug {
		user.ID = 1
		user.Permission = P_ADMIN + P_OPERATOR
		return nil
	}

	// get id
	id, err := GetUserID(c)
	if err != nil {
		return err
	}
	user.ID = id

	// parse JWT
	tokenString := c.Get("Authorization")
	if !user.parseJWT(tokenString) {
		return utils.Unauthorized("Invalid JWT Token")
	}

	for _, v := range user.Roles {
		if v == "admin" {
			user.Permission |= P_ADMIN
		} else if v == "operator" {
			user.Permission |= P_OPERATOR
		} else if strings.HasPrefix(v, "ban_treehole") {
			banDivisionID, err := strconv.Atoi(v[13:]) // "ban_treehole_{divisionID}"
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

// GetAndCheckPermission gets userInfo and check user permission
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

// CheckPermission checks user permission
//
// Example:
//
//  CheckPermission(P_ADMIN)
//  CheckPermission(P_ADMIN | P_OPERATOR)
//
func (user *User) CheckPermission(t PermissionType) bool {
	return user.Permission&t != 0
}
