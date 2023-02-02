// Package penalty is deprecated! Please use APIs in auth.
package penalty

import (
	"github.com/gofiber/fiber/v2"
	"time"
	. "treehole_next/models"
	. "treehole_next/utils"
)

type PostBody struct {
	PenaltyLevel *int `json:"penalty_level" validate:"omitempty"` // low priority, deprecated
	Days         *int `json:"days" validate:"omitempty,min=1"`    // high priority
	DivisionID   int  `json:"division_id" validate:"required,min=1"`
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
	body, err := ValidateBody[PostBody](c)
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
		return Forbidden()
	}

	var floor Floor
	err = DB.Take(&floor, floorID).Error
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

	punishment := Punishment{
		UserID:     floor.UserID,
		MadeBy:     user.ID,
		FloorID:    floor.ID,
		DivisionID: body.DivisionID,
		Duration:   time.Duration(days) * 24 * time.Hour,
		Reason:     "",
	}
	user, err = punishment.Create()
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func RegisterRoutes(app fiber.Router) {
	app.Post("/penalty/:id", BanUser)
}
