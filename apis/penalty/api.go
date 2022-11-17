// Package penalty is deprecated! Please use APIs in auth.
package penalty

import (
	"bytes"
	"encoding/json"
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
	var user User
	err = user.GetUser(c)
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

	banUser(body.DivisionID, days, user.ID, floor.UserID)

	userData := getUser(floor.UserID)

	return c.Send(userData)
}

var client = http.Client{Timeout: time.Second * 10}

func banUser(divisionID int, days int, fromUserID int, toUserID int) {
	data := map[string]any{
		"name":   fmt.Sprintf("ban_treehole_%d", divisionID),
		"days":   days,
		"reason": "ban user",
	}
	dataBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/users/%d/permissions", config.Config.AuthUrl, toUserID),
		bytes.NewBuffer(dataBytes),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Consumer-Username", strconv.Itoa(fromUserID))
	res, err := client.Do(req)
	if err != nil {
		Logger.Error(
			"auth add permission error, request error",
			zap.Int("user id", toUserID),
		)
	}
	if res.StatusCode != 200 {
		Logger.Error(
			"auth add permission error, response error",
			zap.Int("user id", toUserID),
		)
	}
}

func getUser(userID int) []byte {
	response, err := client.Get(fmt.Sprintf("%s/users/%d", config.Config.AuthUrl, userID))
	if err != nil {
		Logger.Error(
			"auth get user error, request error",
			zap.Int("user id", userID),
		)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Logger.Error("close body error", zap.Error(err))
		}
	}(response.Body)

	data, err := io.ReadAll(response.Body)
	if err != nil {
		Logger.Error(
			"auth get user error, read body error",
			zap.Int("user id", userID),
		)
	}

	if response.StatusCode != 200 {
		Logger.Error(
			"auth get user error, response error",
			zap.Int("user id", userID),
			zap.String("response body", string(data)),
		)
	}

	return data
}

func RegisterRoutes(app fiber.Router) {
	app.Post("/penalty/:id", BanUser)
}
