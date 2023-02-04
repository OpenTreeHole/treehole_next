package floor

import (
	"fmt"
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
)

// ListFloorsInAHole
//
//	@Summary	List Floors In A Hole
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/holes/{hole_id}/floors [get]
//	@Param		hole_id	path	int			true	"hole id"
//	@Param		object	query	ListModel	false	"query"
//	@Success	200		{array}	Floor
func ListFloorsInAHole(c *fiber.Ctx) error {
	// validate
	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	query, err := ValidateQuery[ListModel](c)
	if err != nil {
		return err
	}

	// get floors
	var floors Floors
	result := DB.Limit(query.Size).Order(query.OrderBy+" "+query.Sort).
		// use ranking field to locate faster
		Where("hole_id = ? and ranking >= ?", holeID, query.Offset).
		Preload("Mention").
		Find(&floors)
	if result.Error != nil {
		return result.Error
	}

	return Serialize(c, &floors)
}

// ListFloorsOld
//
//	@Summary	Old API for Listing Floors
//	@Deprecated
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/floors [get]
//	@Param		object	query	ListOldModel	false	"query"
//	@Success	200		{array}	Floor
func ListFloorsOld(c *fiber.Ctx) error {
	// validate
	query, err := ValidateQuery[ListOldModel](c)
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
		querySet = DB.Limit(query.Size).Where("ranking >= ?", query.Offset)
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
//
//	@Summary	Get A Floor
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/floors/{id} [get]
//	@Param		id	path		int	true	"id"
//	@Success	200	{object}	Floor
//	@Failure	404	{object}	MessageModel
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
//
//	@Summary	Create A Floor
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/holes/{hole_id}/floors [post]
//	@Param		hole_id	path		int			true	"hole id"
//	@Param		json	body		CreateModel	true	"json"
//	@Success	201		{object}	Floor
func CreateFloor(c *fiber.Ctx) error {
	body, err := ValidateBody[CreateModel](c)
	if err != nil {
		return err
	}

	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get divisionID
	var divisionID int
	result := DB.Table("hole").Select("division_id").Where("id = ?", holeID).Take(&divisionID)
	if result.Error != nil {
		return result.Error
	}

	// get user from auth
	user, err := GetUser(c)
	if err != nil {
		if err != nil {
			return err
		}
	}

	// permission
	if user.BanDivision[divisionID] != nil {
		return Forbidden("您没有权限在此板块发言")
	}

	// create floor
	floor := Floor{
		HoleID:     holeID,
		UserID:     user.ID,
		Content:    body.Content,
		ReplyTo:    body.ReplyTo,
		SpecialTag: body.SpecialTag,
		IsMe:       true,
	}
	err = floor.Create(DB)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&floor)
}

// CreateFloorOld
//
//	@Summary	Old API for Creating A Floor
//	@Deprecated
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/floors [post]
//	@Param		json	body		CreateOldModel	true	"json"
//	@Success	201		{object}	CreateOldResponse
func CreateFloorOld(c *fiber.Ctx) error {
	body, err := ValidateBody[CreateOldModel](c)
	if err != nil {
		return err
	}

	// get divisionID
	var hole Hole
	err = DB.Take(&hole, body.HoleID).Error
	if err != nil {
		return err
	}

	// get user from auth
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission
	if user.BanDivision[hole.DivisionID] != nil {
		return Forbidden("您没有权限在此板块发言")
	}

	// create floor
	floor := Floor{
		HoleID:     body.HoleID,
		UserID:     user.ID,
		Content:    body.Content,
		ReplyTo:    body.ReplyTo,
		SpecialTag: body.SpecialTag,
		IsMe:       true,
	}
	err = floor.Create(DB)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&CreateOldResponse{
		Data:    floor,
		Message: "发表成功",
	})
}

// ModifyFloor
//
//	@Summary		Modify A Floor
//	@Description	when both "fold_v2" and "fold" are empty, reset fold; else, "fold_v2" has the priority
//	@Tags			Floor
//	@Produce		application/json
//	@Router			/floors/{id} [put]
//	@Param			id		path		int			true	"id"
//	@Param			json	body		ModifyModel	true	"json"
//	@Success		200		{object}	Floor
//	@Failure		404		{object}	MessageModel
func ModifyFloor(c *fiber.Ctx) error {
	// validate request body
	body, err := ValidateBody[ModifyModel](c)
	if err != nil {
		return err
	}

	if body.DoNothing() {
		return BadRequest("无效请求")
	}

	// parse floor_id
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// find floor user_id
	var floor Floor
	err = DB.Select("id", "user_id", "hole_id").Take(&floor, floorID).Error
	if err != nil {
		return err
	}

	// find division_id
	var divisionID int
	err = DB.Model(&Hole{}).Select("division_id").Where("id = ?", floor.HoleID).Scan(&divisionID).Error
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// check permission
	err = body.CheckPermission(user, floor.UserID, divisionID)
	if err != nil {
		return err
	}

	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&floor).Error
		if err != nil {
			return err
		}

		// partially modify floor
		if body.Content != nil && *body.Content != "" {
			var reason string
			if user.ID == floor.UserID {
				reason = "该内容已被作者修改"
				MyLog("Floor", "Modify", floorID, user.ID, RoleOwner, "content")
			} else if user.IsAdmin {
				reason = "该内容已被管理员修改"
				MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "content")
			} else {
				return Forbidden()
			}
			floor.Modified += 1
			err = floor.Backup(tx, user.ID, reason)
			if err != nil {
				return err
			}
			floor.Content = *body.Content

			// update floor.mention after update floor.content
			err = tx.Where("floor_id = ?", floorID).Delete(&FloorMention{}).Error
			if err != nil {
				return err
			}

			floor.Mention, err = LoadFloorMentions(tx, floor.Content)
			if err != nil {
				return err
			}

			// save floor_mention association
			if len(floor.Mention) > 0 {
				err = tx.Omit("Mention.*", "UpdatedAt").Select("Mention").Save(&floor).Error
				if err != nil {
					return err
				}
			}
		}

		if body.Fold != nil {
			if *body.Fold != "" {
				floor.Fold = *body.Fold
				MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "fold")
			} else {
				floor.Fold = ""
				MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "fold reset")
			}
		} else if body.FoldFrontend != nil {
			if len(body.FoldFrontend) != 0 {
				floor.Fold = body.FoldFrontend[0]
				MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "fold")
			} else {
				floor.Fold = ""
				MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "fold reset")
			}
		}

		if body.SpecialTag != nil && *body.SpecialTag != "" {
			floor.SpecialTag = *body.SpecialTag
			MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "specialTag to: ", fmt.Sprintf("%v", floor.SpecialTag))
		}

		if body.Like != nil {
			if *body.Like == "add" {
				err = floor.ModifyLike(tx, user.ID, 1)
			} else if *body.Like == "cancel" {
				err = floor.ModifyLike(tx, user.ID, 0)
			}
			if err != nil {
				return err
			}
		}

		// save all maybe-modified fields above
		// including Like when Like == 0
		return tx.Model(&floor).
			Select("Content", "Fold", "SpecialTag", "Like", "DisLike", "Modified").
			Updates(&floor).Error
	})

	// SendModify only when operator or admin modify content or fold
	if ((body.Content != nil && *body.Content != "") ||
		body.Fold != nil || body.FoldFrontend != nil) &&
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
//
//	@Summary	Modify A Floor's like
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/floors/{id}/like/{like} [post]
//	@Param		id		path		int	true	"id"
//	@Param		like	path		int	true	"1 is like, 0 is reset, -1 is dislike"
//	@Success	200		{object}	Floor
//	@Failure	404		{object}	MessageModel
func ModifyFloorLike(c *fiber.Ctx) error {
	// validate like option
	likeOption, err := c.ParamsInt("like")
	if err != nil {
		return err
	}

	// validate like option
	if likeOption > 1 || likeOption < -1 {
		return BadRequest("like option must be -1, 0 or 1")
	}

	// parse floor_id
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	var floor Floor
	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&floor, floorID)
		if result.Error != nil {
			return result.Error
		}

		// modify like
		err = floor.ModifyLike(tx, userID, int8(likeOption))
		if err != nil {
			return err
		}

		// save like only
		return tx.Model(&floor).Select("Like", "Dislike").Updates(&floor).Error
	})
	if err != nil {
		return err
	}

	return Serialize(c, &floor)
}

// DeleteFloor
//
//	@Summary	Delete A Floor
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/floors/{id} [delete]
//	@Param		id		path		int			true	"id"
//	@Param		json	body		DeleteModel	true	"json"
//	@Success	200		{object}	Floor
//	@Failure	404		{object}	MessageModel
func DeleteFloor(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[DeleteModel](c)
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

	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}

	// permission
	if !(user.ID == floor.UserID || user.IsAdmin) {
		return Forbidden()
	}

	err = floor.Backup(DB, user.ID, body.Reason)
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

	return Serialize(c, &floor)
}

// GetFloorHistory
//
//	@Summary	Get A Floor's History, admin only
//	@Tags		Floor
//	@Produce	application/json
//	@Router		/floors/{id}/history [get]
//	@Param		id	path		int	true	"id"
//	@Success	200	{array}		FloorHistory
//	@Failure	404	{object}	MessageModel
func GetFloorHistory(c *fiber.Ctx) error {
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

	var histories []FloorHistory
	result := DB.Where("floor_id = ?", floorID).Find(&histories)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&histories)
}

// RestoreFloor
//
//	@Summary		Restore A Floor, admin only
//	@Description	Restore A Floor From A History Version
//	@Tags			Floor
//	@Router			/floors/{id}/restore/{floor_history_id} [post]
//	@Param			id					path		int				true	"id"
//	@Param			floor_history_id	path		int				true	"floor_history_id"
//	@Param			json				body		RestoreModel	true	"json"
//	@Success		200					{object}	Floor
//	@Failure		404					{object}	MessageModel
func RestoreFloor(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[RestoreModel](c)
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
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission check
	if !user.IsAdmin {
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
	err = floor.Backup(DB, user.ID, reason)
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
