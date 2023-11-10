package floor

import (
	"fmt"
	"slices"
	"time"

	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"

	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
)

// ListFloorsInAHole
//
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
	err = common.ValidateQuery(c, &query)
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
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	if query.Search != "" {
		return SearchFloorsOld(c, &query)
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
//
// @Summary Create A Floor
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [post]
// @Param hole_id path int true "hole id"
// @Param json body CreateModel true "json"
// @Success 201 {object} Floor
func CreateFloor(c *fiber.Ctx) error {
	var body CreateModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	if len([]rune(body.Content)) > 15000 {
		return common.BadRequest("文本限制 15000 字")
	}

	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get hole to check DivisionID and Locked
	var hole Hole
	err = DB.Take(&hole, holeID).Error
	if err != nil {
		return err
	}

	// get user from auth
	user, err := GetUser(c)
	if err != nil {
		if err != nil {
			return err
		}
	}

	// permission
	if user.BanDivision[hole.DivisionID] != nil {
		return common.Forbidden(user.BanDivisionMessage(hole.DivisionID))
	}
	if hole.Locked && !user.IsAdmin {
		return common.Forbidden("该帖子已被锁定，非管理员禁止发帖")
	}

	// special tag
	if body.SpecialTag != "" && !user.IsAdmin && !slices.Contains(user.SpecialTags, body.SpecialTag) {
		return common.Forbidden("非管理员禁止发含有特殊标签的帖")
	} else if body.SpecialTag == "" && user.DefaultSpecialTag != "" {
		body.SpecialTag = user.DefaultSpecialTag
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
	err = floor.Create(DB, &hole)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&floor)
}

// CreateFloorOld
//
// @Summary Old API for Creating A Floor
// @Deprecated
// @Tags Floor
// @Produce application/json
// @Router /floors [post]
// @Param json body CreateOldModel true "json"
// @Success 201 {object} CreateOldResponse
func CreateFloorOld(c *fiber.Ctx) error {
	var body CreateOldModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	if len([]rune(body.Content)) > 15000 {
		return common.BadRequest("文本限制 15000 字")
	}

	// get hole to check DivisionID and Locked
	var hole Hole
	err = DB.Take(&hole, body.HoleID).Error
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission
	if user.BanDivision[hole.DivisionID] != nil {
		return common.Forbidden(user.BanDivisionMessage(hole.DivisionID))
	}
	if hole.Locked && !user.IsAdmin {
		return common.Forbidden("该帖子已被锁定，非管理员禁止发帖")
	}

	// special tag
	if body.SpecialTag != "" && !user.IsAdmin && !slices.Contains(user.SpecialTags, body.SpecialTag) {
		return common.Forbidden("非管理员禁止发含有特殊标签的帖子")
	} else if body.SpecialTag == "" && user.DefaultSpecialTag != "" {
		body.SpecialTag = user.DefaultSpecialTag
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
	err = floor.Create(DB, &hole)
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
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	if body.DoNothing() {
		return common.BadRequest("无效请求")
	}

	if body.Content != nil && len([]rune(*body.Content)) > 15000 {
		return common.BadRequest("文本限制 15000 字")
	}

	// parse floor_id
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
	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		// load floor, lock for update
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&floor, floorID).Error
		if err != nil {
			return err
		}

		// find hole
		var hole Hole
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&hole, floor.HoleID).Error
		if err != nil {
			return err
		}

		// check permission
		err = body.CheckPermission(user, &floor, &hole)
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
				return common.Forbidden()
			}
			floor.Modified += 1
			err = floor.Backup(tx, user.ID, reason)
			if err != nil {
				return err
			}
			floor.Content = *body.Content

			// sensitive check
			err = floor.SensitiveCheck(tx, &hole)
			if err != nil {
				return err
			}

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

			// reindex floor
			go FloorIndex(FloorModel{
				ID:        floor.ID,
				UpdatedAt: time.Now(),
				Content:   floor.Content,
			})
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

		if body.SpecialTag != nil {
			floor.SpecialTag = *body.SpecialTag
			MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "specialTag to: ", fmt.Sprintf("%s", floor.SpecialTag))
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
	if err != nil {
		return err
	}

	// SendModify only when operator or admin modify content or fold
	if ((body.Content != nil && *body.Content != "") ||
		body.Fold != nil || body.FoldFrontend != nil) &&
		user.ID != floor.UserID {
		err = floor.SendModify(DB)
		if err != nil {
			log.Err(err).Str("model", "Notification").Msg("SendModify failed")
			// return err // only for test
		}
	}

	return Serialize(c, &floor)
}

// ModifyFloorLike
//
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

	// validate like option
	if likeOption > 1 || likeOption < -1 {
		return common.BadRequest("like option must be -1, 0 or 1")
	}

	// parse floor_id
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	userID, err := common.GetUserID(c)
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

	var floor Floor
	err = DB.Transaction(func(tx *gorm.DB) error {

		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&floor, floorID)
		if result.Error != nil {
			return result.Error
		}

		// permission
		if !((user.ID == floor.UserID && !floor.Deleted) || user.IsAdmin) {
			return common.Forbidden()
		}

		err = floor.Backup(tx, user.ID, body.Reason)
		if err != nil {
			return err
		}

		floor.Deleted = true
		floor.Content = generateDeleteReason(body.Reason, user.ID == floor.UserID)
		return tx.Save(&floor).Error
	})
	if err != nil {
		return err
	}

	go FloorDelete(floor.ID)

	// log
	if user.ID == floor.UserID {
		MyLog("Floor", "Delete", floorID, user.ID, RoleOperator, "reason: ", body.Reason)
	} else {
		MyLog("Floor", "Delete", floorID, user.ID, RoleOperator, "reason: ", body.Reason)

		// SendModify when admin delete floor
		err = floor.SendModify(DB)
		if err != nil {
			log.Err(err).Str("model", "Notification").Msg("SendModify failed")
			// return err // only for test
		}
	}

	return Serialize(c, &floor)
}

// ListReplyFloors
//
// @Summary List User's Reply Floors
// @Tags Floor
// @Produce application/json
// @Router /users/me/floors [get]
// @Param object query ListModel false "query"
// @Success 200 {array} Floor
// @Failure 404 {object} MessageModel
func ListReplyFloors(c *fiber.Ctx) error {
	//get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	var query ListModel
	err = common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	// get floors
	var floors Floors
	result := DB.
		Limit(query.Size).
		Offset(query.Offset).
		Order(query.OrderBy+" "+query.Sort).
		Where("user_id = ? and ranking <> 0", userID).
		Preload("Mention").
		Find(&floors)
	if result.Error != nil {
		return result.Error
	}

	return Serialize(c, &floors)
}

// GetFloorHistory
//
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
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission
	if !user.IsAdmin {
		return common.Forbidden()
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
	err := common.ValidateBody(c, &body)
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
		return common.Forbidden()
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
		return common.BadRequest(fmt.Sprintf("%v 不是 #%v 的历史版本", floorHistoryID, floorID))
	}
	reason := body.Reason
	err = floor.Backup(DB, user.ID, reason)
	if err != nil {
		return err
	}
	floor.Deleted = false
	floor.Content = floorHistory.Content
	DB.Save(&floor)

	go FloorIndex(FloorModel{
		ID:        floor.ID,
		UpdatedAt: time.Now(),
		Content:   floor.Content,
	})

	// log
	MyLog("Floor", "Restore", floorID, user.ID, RoleAdmin, reason)
	return Serialize(c, &floor)
}

// GetPunishmentHistory
//
// @Summary Get A Floor's Punishment History, admin only
// @Tags Floor
// @Produce application/json
// @Router /floors/{id}/punishment [get]
// @Param id path int true "id"
// @Success 200 {array} string
// @Failure 404 {object} MessageModel
func GetPunishmentHistory(c *fiber.Ctx) error {
	floorID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission, admin only
	if !user.IsAdmin {
		return common.Forbidden()
	}

	// get floor userID
	var floor Floor
	result := DB.First(&floor, floorID)
	if result.Error != nil {
		return result.Error
	}
	userID := floor.UserID

	// search DB for user punishment history
	punishments := make([]string, 0, 10)
	err = DB.Raw(
		`SELECT f.content 
FROM floor f
WHERE f.id IN (
	SELECT distinct floor.id
	FROM floor
	JOIN floor_history ON floor.id = floor_history.floor_id 
	WHERE floor.user_id <> floor_history.user_id 
	AND floor.user_id = ? 
	AND floor.deleted = true
)`, userID).Scan(&punishments).Error
	if err != nil {
		return err
	}
	return c.JSON(punishments)
}

// ListSensitiveFloors
//
// @Summary List sensitive floors, admin only
// @Tags Floor
// @Produce application/json
// @Router /floors/_sensitive [get]
// @Param object query SensitiveFloorRequest false "query"
// @Success 200 {array} SensitiveFloorResponse
// @Failure 404 {object} MessageModel
func ListSensitiveFloors(c *fiber.Ctx) (err error) {
	// validate query
	var query SensitiveFloorRequest
	err = common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission, admin only
	if !user.IsAdmin {
		return common.Forbidden()
	}

	// get floors
	var floors Floors
	querySet := DB
	if query.All == true {
		querySet = querySet.Where("is_sensitive = true")
	} else {
		if query.Open == true {
			querySet = querySet.Where("is_sensitive = true and is_actual_sensitive IS NULL")
		} else {
			querySet = querySet.Where("is_sensitive = true and is_actual_sensitive IS NOT NULL")
		}
	}

	result := querySet.
		Order("updated_at desc").
		Find(&floors)
	if result.Error != nil {
		return result.Error
	}

	var responses = make([]SensitiveFloorResponse, len(floors))
	for i := range responses {
		responses[i].FromModel(floors[i])
	}

	return c.JSON(responses)
}

// ModifyFloorSensitive
//
// @Summary Modify A Floor's actual_sensitive, admin only
// @Tags Floor
// @Produce application/json
// @Router /floors/{id}/_sensitive [put]
// @Param id path int true "id"
// @Param json body ModifySensitiveFloorRequest true "json"
// @Success 200 {object} Floor
// @Failure 404 {object} MessageModel
func ModifyFloorSensitive(c *fiber.Ctx) (err error) {
	// validate body
	var body ModifySensitiveFloorRequest
	err = common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// parse floor_id
	floorID, err := c.ParamsInt("id")
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
		return common.Forbidden()
	}

	var floor Floor
	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&floor, floorID).Error
		if err != nil {
			return err
		}

		// modify actual_sensitive
		floor.IsActualSensitive = &body.IsActualSensitive
		MyLog("Floor", "Modify", floorID, user.ID, RoleAdmin, "actual_sensitive to: ", fmt.Sprintf("%v", body.IsActualSensitive))

		// save actual_sensitive only
		return tx.Model(&floor).Select("IsActualSensitive").UpdateColumns(&floor).Error
	})
	if err != nil {
		return err
	}

	// clear cache
	err = DeleteCache(fmt.Sprintf("hole_%v", floor.HoleID))
	if err != nil {
		return err
	}

	if *floor.IsActualSensitive == false {
		go FloorIndex(FloorModel{
			ID:        floor.ID,
			UpdatedAt: floor.UpdatedAt,
			Content:   floor.Content,
		})
	}

	return Serialize(c, &floor)
}
