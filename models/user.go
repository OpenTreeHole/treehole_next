package models

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	BaseModel
	ID          int                    `json:"id" gorm:"primarykey"`
	Roles       []string               `json:"roles" gorm:"-:all"`
	BanDivision map[int]bool           `json:"-" gorm:"-:all"`
	Nickname    string                 `json:"nickname" gorm:"-:all"`
	Config      map[string]interface{} `json:"config" gorm:"-:all"`
	Permission  PermissionType         `json:"permission" gorm:"-:all"`
}

type UserFavorites struct {
	UserID int `json:"user_id" gorm:"primarykey"`
	HoleID int `json:"hole_id" gorm:"primarykey"`
}

// PermissionType enum
//goland:noinspection GoSnakeCaseUsage
const (
	P_ADMIN = 1 << iota
	P_OPERATOR
)

type PermissionType int

// parseJWT extracts and parse token
func (user *User) parseJWT(token string) error {
	if len(token) < 7 {
		return errors.New("bearer token required")
	}

	payloads := strings.SplitN(token[7:], ".", 3) // extract "Bearer "
	if len(payloads) < 3 {
		return errors.New("jwt token required")
	}

	// jwt encoding ignores padding, so RawStdEncoding should be used instead of StdEncoding
	payloadBytes, err := base64.RawStdEncoding.DecodeString(payloads[1]) // the middle one is payload
	if err != nil {
		return err
	}

	err = json.Unmarshal(payloadBytes, user)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) GetUser(c *fiber.Ctx) error {
	if config.Debug {
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
	err = user.parseJWT(tokenString)
	if err != nil {
		return utils.Unauthorized(err.Error())
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
	if config.Debug {
		return 1, nil
	}

	id, err := strconv.Atoi(c.Get("X-Consumer-Username"))
	if err != nil {
		return 0, utils.Unauthorized("Unauthorized")
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

func UserCreateFavourite(c *fiber.Ctx, clear bool, userID int, holeIDs []int) error {
	if clear {
		DB.Exec("DELETE FROM user_favorites WHERE user_id = ?", userID)
	}
	var builder strings.Builder
	if config.Debug {
		builder.WriteString("INSERT INTO")
	} else {
		builder.WriteString("INSERT IGNORE INTO")
	}
	builder.WriteString(" user_favorites (user_id, hole_id) VALUES ")
	for i, holeID := range holeIDs {
		builder.WriteString(fmt.Sprintf("(%d, %d)", userID, holeID))
		if i != len(holeIDs)-1 {
			builder.WriteString(", ")
		}
	}
	if config.Debug {
		builder.WriteString(" ON CONFLICT DO NOTHING")
	}
	result := DB.Exec(builder.String())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		c.Status(201)
	}
	return nil
}

func UserDeleteFavorite(userID int, holeIDs []int) error {
	sql := "DELETE FROM user_favorites WHERE user_id = ? AND hole_id IN ?"
	result := DB.Exec(sql, userID, holeIDs)
	return result.Error
}
