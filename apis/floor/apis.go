package floor

import (
	"fmt"
	"gorm.io/gorm"
	. "treehole_next/models"
	. "treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"
)

// ListFloorsInAHole
// @Summary List Floors In A Hole
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [get]
// @Param hole_id path int true "hole id"
// @Param object query ListModel false "query"
// @Success 200 {array} Floor
func ListFloorsInAHole(c *fiber.Ctx) error {
	// validate
	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var query ListModel
	err = ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	// get floors
	var floors Floors
	result := query.BaseQuery().
		Where("hole_id = ?", holeID).
		Preload("Mention").
		Find(&floors)
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
// @Param object query ListOldModel false "query"
// @Success 200 {array} Floor
func ListFloorsOld(c *fiber.Ctx) error {
	// validate
	var query ListOldModel

	// parse and validate query manually
	err := c.QueryParser(&query)
	if err != nil {
		return err
	}
	err = Validate(&query)
	if err != nil {
		return err
	}

	if query.Search != "" {
		return SearchFloorsOld(c, query)
	}

	var querySet *gorm.DB
	if query.Size == 0 && query.Offset == 0 {
		querySet = DB
	} else {
		if query.Size == 0 {
			query.Size = 30
		}
		querySet = query.BaseQuery()
	}

	// get floors
	floors := Floors{}
	result := querySet.
		Where("hole_id = ?", query.HoleID).
		Preload("Mention").
		Find(&floors)
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
// @Failure 404 {object} MessageModel
func GetFloor(c *fiber.Ctx) error {
	// validate floor id
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get floor
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

	// create floor
	floor := Floor{
		HoleID:     holeID,
		Content:    body.Content,
		ReplyTo:    body.ReplyTo,
		SpecialTag: body.SpecialTag,
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
// @Success 201 {object} CreateOldResponse
func CreateFloorOld(c *fiber.Ctx) error {
	var body CreateOldModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// create floor
	floor := Floor{
		HoleID:     body.HoleID,
		Content:    body.Content,
		ReplyTo:    body.ReplyTo,
		SpecialTag: body.SpecialTag,
	}
	err = floor.Create(c)
	if err != nil {
		return err
	}

	err = floor.Preprocess(c)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&CreateOldResponse{
		Data:    floor,
		Message: "发表成功",
	})
}

// ModifyFloor
// @Summary Modify A Floor
// @Description when both "fold_v2" and "fold" are empty, reset fold; else, "fold_v2" has the priority
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
			MyLog("Floor", "Modify", floorID, user.ID, RoleOwner, "content")
		} else if perm.CheckPermission(user, perm.Admin) {
			reason = "该内容已被管理员修改"
			MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "content")
		} else {
			return Forbidden()
		}
		err = floor.Backup(c, reason)
		if err != nil {
			return err
		}
		floor.Content = body.Content

		// update floor_mention after update floor.content
		err = floor.SetMention(DB, true)
		if err != nil {
			return err
		}
	}

	if body.Fold == "" && body.FoldFrontend != nil {
		if !perm.CheckPermission(user, perm.Admin) {
			return Forbidden()
		}
		if len(body.FoldFrontend) == 0 {
			// reset floor.Fold
			floor.Fold = ""
			MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "fold reset")
		} else {
			// set floor.Fold
			body.Fold = body.FoldFrontend[0]
		}
	}

	if body.Fold != "" {
		if !perm.CheckPermission(user, perm.Admin) {
			return Forbidden()
		}
		floor.Fold = body.Fold
		MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "fold")
	}

	if body.SpecialTag != "" {
		// operator can modify specialTag
		if !perm.CheckPermission(user, perm.Admin|perm.Operator) {
			return Forbidden()
		}
		floor.SpecialTag = body.SpecialTag
		MyLog("Floor", "Modify", floorID, user.ID, RoleOperator, "specialTag to: ", fmt.Sprintf("%v", floor.SpecialTag))
	}

	if body.Like == "add" {
		err = floor.ModifyLike(c, 1)
	} else if body.Like == "cancel" {
		err = floor.ModifyLike(c, 0)
	}
	if err != nil {
		return err
	}

	// save all fields except Mention
	// including Like when Like == 0
	DB.Model(&floor).Select("*").Omit("Mention").Updates(&floor)

	// SendModify only when operator or admin modify content or fold
	if (body.Content != "" ||
		body.Fold != "") &&
		user.ID != floor.UserID {
		err = floor.SendModify(DB)
		if err != nil {
			Logger.Error("[notification] SendModify failed: " + err.Error())
			// return err // only for test
		}
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
	// validate body
	var body DeleteModel
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

	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}

	// permission
	if !(user.ID == floor.UserID || perm.CheckPermission(user, perm.Admin)) {
		return Forbidden()
	}

	err = floor.Backup(c, body.Reason)
	if err != nil {
		return err
	}

	floor.Deleted = true
	floor.Content = generateDeleteReason(body.Reason, user.ID == floor.UserID)
	DB.Save(&floor)

	// log
	if user.ID == floor.UserID {
		MyLog("Floor", "Delete", floorID, user.ID, RoleOperator, "reason: ", body.Reason)
	} else {
		MyLog("Floor", "Delete", floorID, user.ID, RoleOperator, "reason: ", body.Reason)

		// SendModify when admin delete floor
		err = floor.SendModify(DB)
		if err != nil {
			Logger.Error("[notification] SendModify failed: " + err.Error())
			// return err // only for test
		}
	}

	if user.ID != floor.UserID {

	}

	return Serialize(c, &floor)
}

// GetFloorHistory
// @Summary Get A Floor's History, admin only
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

	var histories []FloorHistory
	result := DB.Where("floor_id = ?", floorID).Find(&histories)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&histories)
}

// RestoreFloor
// @Summary Restore A Floor, admin only
// @Description Restore A Floor From A History Version
// @Tags Floor
// @Router /floors/{id}/restore/{floor_history_id} [post]
// @Param id path int true "id"
// @Param floor_history_id path int true "floor_history_id"
// @Param json body RestoreModel true "json"
// @Success 200 {object} Floor
// @Failure 404 {object} MessageModel
func RestoreFloor(c *fiber.Ctx) error {
	// validate body
	var body RestoreModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get id
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	floorHistoryID, err := c.ParamsInt("floor_history_id")
	if err != nil {
		return err
	}

	// get user
	var user User
	err = user.GetUser(c)
	if err != nil {
		return err
	}

	// permission check
	if !perm.CheckPermission(user, perm.Admin) {
		return Forbidden()
	}

	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}
	var floorHistory FloorHistory
	result = DB.First(&floorHistory, floorHistoryID)
	if result.Error != nil {
		return result.Error
	}
	if floorHistory.FloorID != floorID {
		return BadRequest(fmt.Sprintf("%v 不是 #%v 的历史版本", floorHistoryID, floorID))
	}
	reason := body.Reason
	err = floor.Backup(c, reason)
	if err != nil {
		return err
	}
	floor.Deleted = false
	floor.Content = floorHistory.Content
	DB.Save(&floor)

	// log
	MyLog("Floor", "Restore", floorID, user.ID, RoleAdmin, reason)
	return Serialize(c, &floor)
}
