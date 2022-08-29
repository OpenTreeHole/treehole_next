// Package penalty is deprecated! Please use APIs in auth.
package penalty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"

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

	makeRequest(body.DivisionID, days, floor.UserID)

	return c.Status(200).JSON(nil)
}

var client = http.Client{Timeout: time.Second * 10}

func makeRequest(divisionID int, days int, userID int) {
	data := map[string]any{
		"name":   fmt.Sprintf("ban_treehole_%d", divisionID),
		"days":   days,
		"reason": "ban user",
	}
	dataBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/users/%d/permissions", config.Config.AuthUrl, userID),
		bytes.NewBuffer(dataBytes),
	)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		Logger.Error(
			"auth add permission error, request error",
			zap.Int("user id", userID),
		)
	}
	if res.StatusCode != 200 {
		Logger.Error(
			"auth add permission error, response error",
			zap.Int("user id", userID),
		)
	}
}

func RegisterRoutes(app fiber.Router) {
	app.Post("/penalty/:id", BanUser)
}
