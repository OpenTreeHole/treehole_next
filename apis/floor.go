package apis

import (
	"encoding/json"
	. "treehole_next/config"
	. "treehole_next/models"
	"treehole_next/schemas"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ListFloorsInAHole
// @Summary List Floors In A Hole
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [get]
// @Param hole_id path int true "hole id"
// @Param object query schemas.Query false "query"
// @Success 200 {array} Floor
func ListFloorsInAHole(c *fiber.Ctx) error {
	var query schemas.Query
	err := c.QueryParser(&query)
	if err != nil {
		return err
	}
	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if query.Size == 0 {
		query.Size = Config.Size
	} else if query.Size > Config.MaxSize {
		query.Size = Config.MaxSize
	}
	if query.OrderBy == "" {
		query.OrderBy = "id"
	}

	var floors Floors
	result := Floor{}.MakeQuerySet(
		query.Size, query.Offset, holeID,
		query.OrderBy, query.Desc,
	).Preload("Mention").Find(&floors)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &floors)
}

// ListFloorsOld
// @Summary Old API for Listing Floors
// @Deprecated
// @Tags Floor
// @Produce application/json
// @Router /floors [get]
// @Param object query schemas.ListFloorOld false "query"
// @Success 200 {array} Floor
func ListFloorsOld(c *fiber.Ctx) error {
	var query schemas.ListFloorOld
	err := c.QueryParser(&query)
	if err != nil {
		return err
	}
	var floors Floors
	result := Floor{}.MakeQuerySet(
		query.Size, query.Offset, query.HoleID,
		"id", false,
	).Preload("Mention").Find(&floors)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &floors)
}

// GetFloor
// @Summary Get A Floor
// @Tags Floor
// @Produce application/json
// @Router /floors/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Floor
// @Failure 404 {object} schemas.MessageModel
func GetFloor(c *fiber.Ctx) error {
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var floor Floor
	result := DB.Preload("Mention").First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &floor)
}

// CreateFloor
// @Summary Create A Floor
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [post]
// @Param hole_id path int true "hole id"
// @Param json body schemas.CreateFloor true "json"
// @Success 201 {object} Floor
func CreateFloor(c *fiber.Ctx) error {
	var body schemas.CreateFloor
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	floor := Floor{
		HoleID:  holeID,
		Content: body.Content,
		ReplyTo: body.ReplyTo,
	}
	err = floor.Create(c)
	if err != nil {
		return err
	}

	err = floor.LoadMention()
	if err != nil {
		return err
	}

	return Serialize(c.Status(201), &floor)
}

// CreateFloorOld
// @Summary Old API for Creating A Floor
// @Deprecated
// @Tags Floor
// @Produce application/json
// @Router /floors [post]
// @Param json body schemas.CreateFloorOld true "json"
// @Success 201 {object} Floor
func CreateFloorOld(c *fiber.Ctx) error {
	var body schemas.CreateFloorOld
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	floor := Floor{
		HoleID:  body.HoleID,
		Content: body.Content,
		ReplyTo: body.ReplyTo,
	}
	err = floor.Create(c)
	if err != nil {
		return err
	}

	err = floor.LoadMention()
	if err != nil {
		return err
	}

	return Serialize(c.Status(201), &floor)
}

// ModifyFloor
// @Summary Modify A Floor
// @Tags Floor
// @Produce application/json
// @Router /floors/{id} [put]
// @Param id path int true "id"
// @Param json body schemas.ModifyFloor true "json"
// @Success 200 {object} Floor
// @Failure 404 {object} schemas.MessageModel
func ModifyFloor(c *fiber.Ctx) error {
	var body schemas.ModifyFloor
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var bodyMap StringMap
	err = json.Unmarshal(c.Body(), &bodyMap)
	if err != nil {
		return err
	}

	// find floor
	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}

	// partially modify floor
	if body.Content != "" {
		floor.Content = body.Content
	}
	if body.SpecialTag != "" {
		floor.SpecialTag = body.SpecialTag
	}
	if _, ok := bodyMap["like_int"]; ok {
		if body.Like != 0 {
			floor.Like += body.Like
		} else {
			floor.Like = 0
		}
	} else if _, ok := bodyMap["like"]; ok {
		if body.LikeOld == "add" {
			floor.Like += 1
		} else if body.LikeOld == "reset" {
			floor.Like = 0
		}
	}
	if body.Fold != "" {
		floor.Fold = body.Fold
	}
	DB.Save(&floor)
	err = floor.LoadMention()
	if err != nil {
		return err
	}
	return Serialize(c, &floor)
}

// DeleteFloor
// @Summary Delete A Floor
// @Tags Floor
// @Produce application/json
// @Router /floors/{id} [delete]
// @Param id path int true "id"
// @Param json body schemas.DeleteFloor true "json"
// @Success 200 {object} Floor
// @Failure 404 {object} schemas.MessageModel
func DeleteFloor(c *fiber.Ctx) error {
	var body schemas.DeleteFloor
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var floor Floor
	floor.ID = floorID

	var floorHistory FloorHistory
	floorHistory.FloorID = floorID
	floorHistory.Content = floor.Content
	floorHistory.Reason = body.Reason

	err = DB.Transaction(func(tx *gorm.DB) error {
		result := DB.Model(&floor).Select("Deleted").Updates(Floor{Deleted: true})
		if result.Error != nil {
			return result.Error
		}
		result = DB.Create(&floorHistory)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}

	return Serialize(c, &floor)
}
