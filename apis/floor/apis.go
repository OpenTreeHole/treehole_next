package floor

import (
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

// ListFloorsInAHole
// @Summary List Floors In A Hole
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [get]
// @Param hole_id path int true "hole id"
// @Param object query Query false "query"
// @Success 200 {array} Floor
func ListFloorsInAHole(c *fiber.Ctx) error {
	// validate
	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var query Query
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	// get floors
	var floors Floors
	result := BaseQuery(&query).Where("hole_id = ?", holeID).Find(&floors)
	if result.Error != nil {
		return result.Error
	}
	err = floors.LoadDyField(c)
	if err != nil {
		return err
	}

	return Serialize(c, &floors)
}

// ListFloorsOld
// @Summary Old API for Listing Floors
// @Deprecated
// @Tags Floor
// @Produce application/json
// @Router /floors [get]
// @Param object query ListOldModel false "query"
// @Success 200 {array} Floor
func ListFloorsOld(c *fiber.Ctx) error {
	// validate
	var query ListOldModel
	err := ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	// get floors
	var floors Floors
	result := BaseQuery(&Query{
		Size:    query.Size,
		Offset:  query.Offset,
		OrderBy: "storey",
	}).Where("hole_id = ?", query.HoleID).Find(&floors)
	if result.Error != nil {
		return result.Error
	}
	err = floors.LoadDyField(c)
	if err != nil {
		return err
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
// @Failure 404 {object} MessageModel
func GetFloor(c *fiber.Ctx) error {
	// validate floor id
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get floor
	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}
	err = floor.LoadDyField(c)
	if err != nil {
		return err
	}

	return Serialize(c, &floor)
}

// CreateFloor
// @Summary Create A Floor
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [post]
// @Param hole_id path int true "hole id"
// @Param json body CreateModel true "json"
// @Success 201 {object} Floor
func CreateFloor(c *fiber.Ctx) error {
	var body CreateModel
	err := ValidateBody(c, &body)
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
// @Param json body CreateOldModel true "json"
// @Success 201 {object} Floor
func CreateFloorOld(c *fiber.Ctx) error {
	var body CreateOldModel
	err := ValidateBody(c, &body)
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
// @Param json body ModifyModel true "json"
// @Success 200 {object} Floor
// @Failure 404 {object} MessageModel
func ModifyFloor(c *fiber.Ctx) error {
	// validate request body
	var body ModifyModel
	err := ValidateBody(c, &body)
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

		// find mention
		err := floor.LoadMention()
		if err != nil {
			return err
		}
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
		err = floor.ModifyLike(c, 1)
	} else if body.Like == "reset" {
		err = floor.ModifyLike(c, 0)
	}
	if err != nil {
		return err
	}

	DB.Save(&floor)
	err = floor.LoadDyField(c)
	if err != nil {
		return err
	}

	return Serialize(c, &floor)
}

// ModifyFloorLike
// @Summary Modify A Floor's like
// @Tags Floor
// @Produce application/json
// @Router /floors/{id}/like/{like} [post]
// @Param id path int true "id"
// @Param like path int true "1 is like, 0 is reset, -1 is dislike"
// @Success 200 {object} Floor
// @Failure 404 {object} MessageModel
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

	// modify like
	err = floor.ModifyLike(c, int8(likeOption))
	if err != nil {
		return err
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
// @Param json body DeleteModel true "json"
// @Success 200 {object} Floor
// @Failure 404 {object} MessageModel
func DeleteFloor(c *fiber.Ctx) error {
	var body DeleteModel
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

	err = floor.Backup(c, body.Reason)
	if err != nil {
		return err
	}

	floor.Deleted = true
	DB.Save(&floor)

	err = floor.LoadDyField(c)
	if err != nil {
		return err
	}
	return Serialize(c, &floor)
}

// GetFloorHistory
// @Summary Get A Floor's History
// @Tags Floor
// @Produce application/json
// @Router /floors/{id}/history [get]
// @Param id path int true "id"
// @Success 200 {array} FloorHistory
// @Failure 404 {object} MessageModel
func GetFloorHistory(c *fiber.Ctx) error {
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var histories []FloorHistory
	result := DB.Where("floor_id = ?", floorID).Find(&histories)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&histories)
}
