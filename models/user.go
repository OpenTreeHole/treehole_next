package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID          int             `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time       `json:"time_created"`
	UpdatedAt   time.Time       `json:"time_updated"`
	Roles       []string        `json:"roles" gorm:"-:all"`
	BanDivision map[int]bool    `json:"-" gorm:"-:all"`
	Nickname    string          `json:"nickname" gorm:"-:all"`
	Config      map[string]any  `json:"config" gorm:"-:all"`
	Permission  perm.Permission `json:"-" gorm:"-:all"`
}

func (user User) GetID() int {
	return user.ID
}

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

func GetUser(c *fiber.Ctx) (*User, error) {
	user := &User{
		Roles:       make([]string, 0, 10),
		BanDivision: make(map[int]bool),
		Config:      make(map[string]any),
		Permission:  0,
	}
	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		user.ID = 1
		user.Permission = perm.Admin + perm.Operator
		return user, nil
	}

	// get id
	id, err := GetUserID(c)
	if err != nil {
		return nil, err
	}
	user.ID = id

	// parse JWT
	tokenString := c.Get("Authorization")
	if tokenString == "" { // token can be in either header or cookie
		tokenString = c.Cookies("access")
	}
	err = user.parseJWT(tokenString)
	if err != nil {
		return nil, utils.Unauthorized(err.Error())
	}

	err = user.parsePermission()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (user *User) parsePermission() error {
	for _, v := range user.Roles {
		if v == "admin" {
			user.Permission |= perm.Admin
		} else if v == "operator" {
			user.Permission |= perm.Operator
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

func GetUserFromAuth(c *fiber.Ctx) (*User, error) {
	user := &User{
		Roles:       make([]string, 0, 10),
		BanDivision: make(map[int]bool),
		Config:      make(map[string]any),
		Permission:  0,
	}

	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		user.ID = 1
		user.Permission = perm.Admin + perm.Operator
		return user, nil
	}

	userID, err := GetUserID(c)
	if err != nil {
		return nil, err
	}

	// make request
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/users/%d", config.Config.AuthUrl, userID),
		bytes.NewBuffer(make([]byte, 0, 10)),
	)
	if err != nil {
		utils.Logger.Error("request make error", zap.Error(err))
		return nil, err
	}

	// add headers
	req.Header.Add("X-Consumer-Username", strconv.Itoa(userID))
	rsp, err := client.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Logger.Error("close rsp body error")
		}
	}(rsp.Body)

	if err != nil {
		utils.Logger.Error(
			"auth get user request error",
			zap.Int("user id", userID),
		)
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, errors.New("auth get user error, rsp error")
	}

	userInfo, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(userInfo, user)
	if err != nil {
		return nil, err
	}

	err = user.parsePermission()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserID(c *fiber.Ctx) (int, error) {
	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		return 1, nil
	}

	id, err := strconv.Atoi(c.Get("X-Consumer-Username"))
	if err != nil {
		return 0, utils.Unauthorized("Unauthorized")
	}

	return id, nil
}

func (user *User) GetPermission() perm.Permission {
	return user.Permission
}
