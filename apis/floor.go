package apis

import (
	. "treehole_next/config"
	. "treehole_next/models"
	"treehole_next/schemas"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
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
	).Find(&floors)
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
	).Find(&floors)
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
	// validate request body
	var body schemas.ModifyFloor
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	// find floor
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}

	// get user
	var user User
	err = user.GetUser(c)
	if err != nil {
		return err
	}

	// partially modify floor
	if body.Content != "" {
		var reason string
		if user.ID == floor.UserID {
			reason = "该内容已被作者修改"
		} else if user.IsAdmin {
			reason = "该内容已被管理员修改"
		} else {
			return Forbidden()
		}
		err = floor.Backup(c, reason)
		if err != nil {
			return err
		}
		floor.Content = body.Content
	}

	if body.Fold != "" {
		if !user.IsAdmin {
			return Forbidden()
		}
		floor.Fold = body.Fold
	}

	if body.SpecialTag != "" {
		if !user.IsAdmin {
			return Forbidden()
		}
		floor.SpecialTag = body.SpecialTag
	}

	if body.Like == "add" {
		floor.ModifyLike(c, true, false)
	} else if body.Like == "reset" {
		floor.ModifyLike(c, false, true)
	}

	DB.Save(&floor)

	return Serialize(c, &floor)
}

// ModifyFloorLike
// @Summary Modify A Floor's like
// @Tags Floor
// @Produce application/json
// @Router /floors/{id}/like/{like} [post]
// @Param id path int true "id"
// @Param like path int true "like"
// @Success 200 {object} Floor
// @Failure 404 {object} schemas.MessageModel
func ModifyFloorLike(c *fiber.Ctx) error {
	// validate like option
	likeOption, err := c.ParamsInt("like")
	if err != nil {
		return err
	}

	// find floor
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}

	// modify LikeMapping and Floor->Like
	if likeOption > 0 {
		floor.ModifyLike(c, true, false)
	} else if likeOption < 0 {
		floor.ModifyLike(c, false, false)
	} else {
		floor.ModifyLike(c, false, true)
	}

	DB.Save(&floor)

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
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}

	err = floor.Backup(c, body.Reason)
	if err != nil {
		return err
	}

	floor.Deleted = true
	DB.Save(&floor)

	return Serialize(c, &floor)
}
