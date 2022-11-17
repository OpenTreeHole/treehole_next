// Package penalty is deprecated! Please use APIs in auth.
package penalty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PostBody struct {
	PenaltyLevel int `json:"penalty_level"`
	DivisionID   int `json:"division_id"`
}

// BanUser
// @Summary [Deprecated] Ban publisher of a floor
// @Deprecated
// @Tags Penalty
// @Produce application/json
// @Router /penalty/{floor_id} [post]
// @Param json body PostBody true "json"
// @Success 201 {object} Hole
func BanUser(c *fiber.Ctx) error {
	// check AuthURL
	if config.Config.AuthUrl == "" {
		return BadRequest("No AuthURL")
	}

	// validate body
	var body PostBody
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission
	if !perm.CheckPermission(user, perm.Admin) {
		return Forbidden()
	}

	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}

	var days int
	switch body.PenaltyLevel {
	case 1:
		days = 1
	case 2:
		days = 5
	case 3:
		days = 999
	default:
		days = 1
	}

	err = banUser(c, body.DivisionID, days, user.ID, floor.UserID)
	if err != nil {
		return err
	}

	userData, err := getUser(floor.UserID)
	if err != nil {
		return err
	}

	return c.SendStream(userData)
}

var client = http.Client{Timeout: time.Second * 10}

func banUser(c *fiber.Ctx, divisionID, days, fromUserID, toUserID int) error {
	data := map[string]any{
		"name":   fmt.Sprintf("ban_treehole_%d", divisionID),
		"days":   days,
		"reason": "ban user",
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// make request
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/users/%d/permissions", config.Config.AuthUrl, toUserID),
		bytes.NewBuffer(dataBytes),
	)
	if err != nil {
		Logger.Error("req make err", zap.Error(err))
		return err
	}

	// add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Consumer-Username", strconv.Itoa(fromUserID))
	res, err := client.Do(req)

	defer func(body io.Closer) {
		err := body.Close()
		if err != nil {
			Logger.Error("close request body error", zap.Error(err))
		}
	}(res.Body)

	// error handling
	if err != nil {
		Logger.Error(
			"auth add permission error, request error",
			zap.Int("user id", toUserID),
		)
		return err
	}
	if res.StatusCode != 200 {
		Logger.Error(
			"auth add permission error, response error",
			zap.Int("user id", toUserID),
		)
		return fiber.ErrInternalServerError
	}

	return nil
}

func getUser(toUserID int) (io.Reader, error) {
	// make request
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/users/%d", config.Config.AuthUrl, toUserID),
		bytes.NewBuffer(make([]byte, 0, 10)),
	)
	if err != nil {
		Logger.Error("request make error", zap.Error(err))
		return nil, err
	}

	// add headers
	req.Header.Add("X-Consumer-Username", strconv.Itoa(toUserID))
	rsp, err := client.Do(req)

	if err == nil && rsp.StatusCode == 200 {
		// do not close body. io.Reader will send to fiber context
		return rsp.Body, nil
	}

	closeErr := rsp.Body.Close()
	if closeErr != nil {
		Logger.Error("close request body error", zap.Error(closeErr))
	}

	// error handling
	if err != nil {
		Logger.Error(
			"auth get user error, request error",
			zap.Int("user id", toUserID),
		)
		return nil, err
	}

	if rsp.StatusCode != 200 {
		Logger.Error(
			"auth get user error, rsp error",
			zap.Int("user id", toUserID),
		)
		return nil, errors.New("auth get user error, rsp error")
	}

	return nil, fiber.ErrInternalServerError
}

func RegisterRoutes(app fiber.Router) {
	app.Post("/penalty/:id", BanUser)
}
