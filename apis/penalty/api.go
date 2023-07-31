// Package penalty is deprecated! Please use APIs in auth.
package penalty

import (
	"fmt"
	"time"

	"github.com/opentreehole/go-common"

	. "treehole_next/models"

	"github.com/gofiber/fiber/v2"
)

type PostBody struct {
	PenaltyLevel *int   `json:"penalty_level" validate:"omitempty"` // low priority, deprecated
	Days         *int   `json:"days" validate:"omitempty,min=1"`    // high priority
	Reason       string `json:"reason"`                             // optional
}

// BanUser
//
//	@Summary	Ban publisher of a floor
//	@Tags		Penalty
//	@Produce	json
//	@Router		/penalty/{floor_id} [post]
//	@Param		json	body		PostBody	true	"json"
//	@Success	201		{object}	User
func BanUser(c *fiber.Ctx) error {
	// validate body
	var body PostBody
	err := common.ValidateBody(c, &body)
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
	if !user.IsAdmin {
		return common.Forbidden()
	}

	var floor Floor
	err = DB.Take(&floor, floorID).Error
	if err != nil {
		return err
	}

	var hole Hole
	err = DB.Take(&hole, floor.HoleID).Error
	if err != nil {
		return err
	}

	var days int
	if body.Days != nil {
		days = *body.Days
		if days <= 0 {
			days = 1
		}
	} else if body.PenaltyLevel != nil {
		switch *body.PenaltyLevel {
		case 1:
			days = 1
		case 2:
			days = 5
		case 3:
			days = 999
		default:
			days = 1
		}
	}

	duration := time.Duration(days) * 24 * time.Hour

	punishment := Punishment{
		UserID:     floor.UserID,
		MadeBy:     user.ID,
		FloorID:    &floor.ID,
		DivisionID: hole.DivisionID,
		Duration:   &duration,
		Reason:     body.Reason,
	}
	user, err = punishment.Create()
	if err != nil {
		return err
	}

	// construct message for user
	message := Notification{
		Data:       floor,
		Recipients: []int{floor.UserID},
		Description: fmt.Sprintf(
			"分区：%d，时间：%d天，原因：%s",
			hole.DivisionID,
			days,
			body.Reason,
		),
		Title: "您的权限被禁止了",
		Type:  MessageTypePermission,
		URL:   fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	// send
	_, err = message.Send()
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// ListMyPunishments godoc
// @Summary List my punishments
// @Tags Penalty
// @Produce json
// @Router /users/me/punishments [get]
// @Success 200 {array} Punishment
func ListMyPunishments(c *fiber.Ctx) error {
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	punishments, err := listPunishmentsByUserID(userID)
	if err != nil {
		return err
	}

	return c.JSON(punishments)
}

// ListPunishmentsByUserID godoc
// @Summary List punishments by user id
// @Tags Penalty
// @Produce json
// @Router /users/{id}/punishments [get]
// @Param id path int true "User ID"
// @Success 200 {array} Punishment
func ListPunishmentsByUserID(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	currentUser, err := GetUser(c)
	if !currentUser.IsAdmin && currentUser.ID != userID {
		return common.Forbidden()
	}

	punishments, err := listPunishmentsByUserID(userID)
	if err != nil {
		return err
	}

	return c.JSON(punishments)
}

func listPunishmentsByUserID(userID int) ([]Punishment, error) {
	var punishments []Punishment
	err := DB.Where("user_id = ?", userID).Preload("Floor").Find(&punishments).Error
	if err != nil {
		return nil, err
	}

	// remove made_by
	for i := range punishments {
		punishments[i].MadeBy = 0
	}

	return punishments, nil
}

func RegisterRoutes(app fiber.Router) {
	app.Post("/penalty/:id", BanUser)
	app.Get("/users/me/punishments", ListMyPunishments)
	app.Get("/users/:id/punishments", ListPunishmentsByUserID)
}
