package hole

import (
	"fmt"
	"slices"
	"strconv"
	"time"
	"treehole_next/utils/sensitive"

	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	. "treehole_next/models"
	"treehole_next/utils"
	. "treehole_next/utils"
)

// ListHolesByDivision
//
// @Summary List Holes In A Division
// @Tags Hole
// @Produce json
// @Router /divisions/{division_id}/holes [get]
// @Param division_id path int true "division_id"
// @Param object query QueryTime false "query"
// @Success 200 {array} Hole
// @Failure 404 {object} MessageModel
// @Failure 500 {object} MessageModel
func ListHolesByDivision(c *fiber.Ctx) error {
	var query QueryTime
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get holes
	var holes Holes
	querySet, err := holes.MakeQuerySet(query.Offset, query.Size, query.Order, c)
	if err != nil {
		return err
	}
	if id != 0 {
		querySet = querySet.Where("hole.division_id = ?", id)
	}
	querySet.Find(&holes)

	return Serialize(c, &holes)
}

// ListHolesByTag
//
// @Summary List Holes By Tag
// @Tags Hole
// @Produce json
// @Router /tags/{tag_name}/holes [get]
// @Param tag_name path string true "tag_name"
// @Param object query QueryTime false "query"
// @Success 200 {array} Hole
// @Failure 404 {object} MessageModel
func ListHolesByTag(c *fiber.Ctx) error {
	var query QueryTime
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	// get tag
	var tag Tag
	tagName := c.Params("name")
	result := DB.Where("name = ?", tagName).First(&tag)
	if result.Error != nil {
		return result.Error
	}

	// get holes
	var holes Holes
	querySet, err := holes.MakeQuerySet(query.Offset, query.Size, "", c)
	if err != nil {
		return err
	}
	err = querySet.Model(&tag).
		Association("Holes").Find(&holes)
	if err != nil {
		return err
	}

	return Serialize(c, &holes)
}

// ListHolesByMe
//
// @Summary List a Hole Created By User
// @Tags Hole
// @Produce json
// @Router /users/me/holes [get]
// @Param object query QueryTime false "query"
// @Success 200 {array} Hole
func ListHolesByMe(c *fiber.Ctx) error {
	var query QueryTime
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// get holes
	var holes Holes
	querySet, err := holes.MakeQuerySet(query.Offset, query.Size, "", c)
	if err != nil {
		return err
	}
	querySet = querySet.Where("hole.user_id = ?", userID)
	querySet.Find(&holes)

	return Serialize(c, &holes)
}

// ListGoodHoles
//
// @Summary List good holes
// @Tags Hole
// @Produce json
// @Router /holes/_good [get]
// @Param object query QueryTime false "query"
// @Success 200 {array} Hole
func ListGoodHoles(c *fiber.Ctx) error {
	var query QueryTime
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}
	_, err = common.GetUserID(c)
	if err != nil {
		return err
	}

	// get holes
	var holes Holes
	querySet, err := holes.MakeQuerySet(query.Offset, query.Size, query.Order, c)
	if err != nil {
		return err
	}
	querySet = querySet.Where("hole.good = 1")
	err = querySet.Find(&holes).Error
	if err != nil {
		return err
	}

	return Serialize(c, &holes)
}

// ListHolesOld
//
// @Summary Old API for Listing Holes
// @Deprecated
// @Tags Hole
// @Produce json
// @Router /holes [get]
// @Param object query ListOldModel false "query"
// @Success 200 {array} Hole
func ListHolesOld(c *fiber.Ctx) error {
	var query ListOldModel
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	var holes Holes
	querySet, err := holes.MakeQuerySet(query.Offset, query.Size, query.Order, c)
	if err != nil {
		return err
	}
	if query.Tag != "" {
		var tag Tag
		err = DB.Where("name = ?", query.Tag).Find(&tag).Error
		if err != nil {
			return err
		}
		err = querySet.Model(&tag).Order("updated_at desc").Association("Holes").Find(&holes)
		if err != nil {
			return err
		}
	} else if query.DivisionID != 0 {
		querySet.
			Where("hole.division_id = ?", query.DivisionID).
			Find(&holes)
	} else {
		querySet.Find(&holes)
	}

	// only for danxi v1.3.10 old api
	if query.Order == "time_created" || query.Order == "created_at" {
		for i := range holes {
			holes[i].UpdatedAt = holes[i].CreatedAt
		}
	}

	return Serialize(c, &holes)
}

// GetHole
//
// @Summary Get A Hole
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Hole
// @Failure 404 {object} MessageModel
func GetHole(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")

	querySet, err := MakeHoleQuerySet(c)
	if err != nil {
		return err
	}

	// get hole
	var hole Hole
	err = querySet.Take(&hole, id).Error
	if err != nil {
		return err
	}

	return Serialize(c, &hole)
}

// CreateHole
//
// @Summary Create A Hole
// @Description Create a hole, create tags and floor binding to it and set the name mapping
// @Tags Hole
// @Produce application/json
// @Router /divisions/{division_id}/holes [post]
// @Param division_id path int true "division id"
// @Param json body CreateModel true "json"
// @Success 201 {object} Hole
func CreateHole(c *fiber.Ctx) error {
	// validate body
	var body CreateModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	if len([]rune(body.Content)) > 10000 {
		return common.BadRequest("文本限制 10000 字")
	}

	divisionID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get user from auth
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission
	if user.BanDivision[divisionID] != nil {
		return common.Forbidden(user.BanDivisionMessage(divisionID))
	}

	// special tag
	if body.SpecialTag != "" && !user.IsAdmin && !slices.Contains(user.SpecialTags, body.SpecialTag) {
		return common.Forbidden("非管理员禁止发含有特殊标签的洞")
	} else if body.SpecialTag == "" && user.DefaultSpecialTag != "" {
		body.SpecialTag = user.DefaultSpecialTag
	}

	sensitiveResp, err := sensitive.CheckSensitive(sensitive.ParamsForCheck{
		Content:  body.Content,
		Id:       time.Now().UnixNano(),
		TypeName: sensitive.TypeFloor,
	})
	if err != nil {
		return err
	}

	hole := Hole{
		Floors: Floors{{
			UserID:          user.ID,
			Content:         body.Content,
			SpecialTag:      body.SpecialTag,
			IsMe:            true,
			IsSensitive:     !sensitiveResp.Pass,
			SensitiveDetail: sensitiveResp.Detail,
		}},
		UserID:     user.ID,
		DivisionID: divisionID,
	}
	err = hole.Create(DB, user, body.ToName(), c)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&hole)
}

// CreateHoleOld
//
// @Summary Old API for Creating A Hole
// @Deprecated
// @Tags Hole
// @Produce application/json
// @Router /holes [post]
// @Param json body CreateOldModel true "json"
// @Success 201 {object} CreateOldResponse
func CreateHoleOld(c *fiber.Ctx) error {
	// validate body
	var body CreateOldModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	if len([]rune(body.Content)) > 10000 {
		return common.BadRequest("文本限制 10000 字")
	}

	// get user from auth
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission
	if user.BanDivision[body.DivisionID] != nil {
		return common.Forbidden(user.BanDivisionMessage(body.DivisionID))
	}

	// special tag
	if body.SpecialTag != "" && !user.IsAdmin && !slices.Contains(user.SpecialTags, body.SpecialTag) {
		return common.Forbidden("非管理员禁止发含有特殊标签的洞")
	} else if body.SpecialTag == "" && user.DefaultSpecialTag != "" {
		body.SpecialTag = user.DefaultSpecialTag
	}

	sensitiveResp, err := sensitive.CheckSensitive(sensitive.ParamsForCheck{
		Content:  body.Content,
		Id:       time.Now().UnixNano(),
		TypeName: sensitive.TypeFloor,
	})
	if err != nil {
		return err
	}

	// create hole
	hole := Hole{
		Floors: Floors{{
			UserID:          user.ID,
			Content:         body.Content,
			SpecialTag:      body.SpecialTag,
			IsMe:            true,
			IsSensitive:     !sensitiveResp.Pass,
			SensitiveDetail: sensitiveResp.Detail,
		}},
		UserID:     user.ID,
		DivisionID: body.DivisionID,
	}
	err = hole.Create(DB, user, body.ToName(), c)
	if err != nil {
		return err
	}

	err = hole.Preprocess(c)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(&CreateOldResponse{
		Data:    hole,
		Message: "发表成功",
	})
}

// ModifyHole
//
// @Summary Modify A Hole
// @Description Modify a hole, modify tags and set the name mapping
// @Description Only admin can modify division, tags, hidden, lock
// @Description `unhidden` take effect only when hole is hidden and set to true
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [put]
// @Param id path int true "id"
// @Param json body ModifyModel true "json"
// @Success 200 {object} Hole
// @Failure 404 {object} MessageModel
func ModifyHole(c *fiber.Ctx) error {
	// validate
	var body ModifyModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}
	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	if body.DoNothing() {
		return common.BadRequest("无效请求")
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// load hole
	var hole Hole

	changed := false

	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		// lock for update
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&hole, holeID).Error
		if err != nil {
			return err
		}

		// check user permission
		err = body.CheckPermission(user, &hole)
		if err != nil {
			return err
		}

		// modify division
		if body.DivisionID != nil && *body.DivisionID != 0 && *body.DivisionID != hole.DivisionID {
			hole.DivisionID = *body.DivisionID
			changed = true
			// log
			MyLog("Hole", "Modify", holeID, user.ID, RoleAdmin, "DivisionID to: ", strconv.Itoa(hole.DivisionID))
		}

		// modify hidden
		if body.Unhidden != nil && *body.Unhidden && hole.Hidden {
			hole.Hidden = false
			changed = true

			// reindex into Elasticsearch
			var floors Floors
			_ = DB.Where("hole_id = ?", hole.ID).Find(&floors)
			var floorModels []FloorModel
			for _, floor := range floors {
				floorModels = append(floorModels, FloorModel{
					ID:        floor.ID,
					UpdatedAt: floor.UpdatedAt,
					Content:   floor.Content,
				})
			}
			go BulkInsert(floorModels)

			// log
			MyLog("Hole", "Modify", holeID, user.ID, RoleAdmin, "Unhidden: ")
		}

		// modify tags
		if len(body.Tags) != 0 {
			changed = true
			hole.Tags, err = FindOrCreateTags(tx, user, body.ToName())
			if err != nil {
				return err
			}

			// set tag.temperature = tag.temperature - 1
			err = tx.Model(&Tag{}).Where("id in (?)",
				tx.Model(&HoleTag{}).Select("tag_id").Where("hole_id = ?", hole.ID)).
				Update("temperature", gorm.Expr("temperature - 1")).Error
			if err != nil {
				return err
			}

			// delete old hole_tags association
			err = tx.Exec("DELETE FROM hole_tags WHERE hole_id = ?", hole.ID).Error
			if err != nil {
				return err
			}

			// Create hole_tags association only
			err = tx.Omit("Tags.*", "UpdatedAt").Select("Tags").Save(&hole).Error
			if err != nil {
				return err
			}

			// Update tag temperature
			err = tx.Model(&hole.Tags).Update("temperature", gorm.Expr("temperature + 1")).Error
			if err != nil {
				return err
			}

			if user.IsAdmin {
				MyLog("Hole", "Modify", holeID, user.ID, RoleAdmin, "NewTags: ", fmt.Sprintf("%v", body.Tags))
			} else {
				MyLog("Hole", "Modify", holeID, user.ID, RoleOwner, "NewTags: ", fmt.Sprintf("%v", body.Tags))
			}
		}

		// modify lock
		if body.Lock != nil {
			changed = true
			hole.Locked = *body.Lock

			MyLog("Hole", "Modify", holeID, user.ID, RoleAdmin, "Lock: ")
		}

		// save
		if changed {
			err = tx.Model(&hole).
				Omit(clause.Associations, "UpdatedAt").
				Select("DivisionID", "Hidden", "Locked").
				Updates(&hole).Error
			if err != nil {
				return err
			}

			if user.IsAdmin {
				CreateAdminLog(tx, AdminLogTypeHole, user.ID, body)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// update cache
	if changed {
		err = UpdateHoleCache(Holes{&hole})
		if err != nil {
			return err
		}
	}

	return Serialize(c, &hole)
}

// HideHole
//
// @Summary Delete A Hole
// @Description Hide a hole, but visible to admins. This may affect many floors, DO NOT ABUSE!!!
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [delete]
// @Param id path int true "id"
// @Success 204
// @Failure 404 {object} MessageModel
func HideHole(c *fiber.Ctx) error {
	// validate holeID
	holeID, err := c.ParamsInt("id")
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

	var hole Hole
	hole.ID = holeID
	result := DB.Model(&hole).Select("Hidden").Omit("UpdatedAt").Updates(Hole{Hidden: true})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	// log
	MyLog("Hole", "Hide", holeID, user.ID, RoleAdmin)

	// find hole and update cache

	err = DB.Take(&hole).Error
	if err != nil {
		return err
	}

	updateHoles := Holes{&hole}
	err = UpdateHoleCache(updateHoles)
	if err != nil {
		return err
	}

	// delete floors from Elasticsearch
	var floors Floors
	_ = DB.Where("hole_id = ?", hole.ID).Find(&floors)
	go BulkDelete(Models2IDSlice(floors))

	return c.Status(204).JSON(nil)
}

// PatchHole
//
// @Summary Patch A Hole
// @Description Add hole.view
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [patch]
// @Param id path int true "id"
// @Success 204
// @Failure 404 {object} MessageModel
func PatchHole(c *fiber.Ctx) error {
	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	holeViewsChan <- holeID

	return c.Status(204).JSON(nil)
}

// DeleteHole godoc
//
// @Summary Delete A Hole
// @Description Delete a hole, admin only
// @Tags Hole
// @Produce json
// @Router /holes/{id}/_force [delete]
// @Param id path int true "id"
// @Success 204
// @Failure 401 {object} MessageModel "Unauthorized"
// @Failure 404 {object} MessageModel "Not Found"
func DeleteHole(c *fiber.Ctx) error {
	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	user, err := GetUser(c)
	if err != nil {
		return err
	}

	var userType Role = RoleOwner

	if user.IsAdmin {
		userType = RoleAdmin
	}

	var hole Hole
	err = DB.Take(&hole, holeID).Error
	if err != nil {
		return err
	}

	if hole.UserID != user.UserID && !user.IsAdmin {
		return common.Forbidden()
	}

	result := DB.Where("hole_id = ? ", holeID).Delete(&hole)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	MyLog("Hole", "Delete", holeID, user.ID, userType)

	err = utils.DeleteCache(hole.CacheName())
	if err != nil {
		log.Err(err).Msg("DeleteHole: delete cache")
	}

	// delete floors from Elasticsearch
	var floors Floors
	err = DB.Where("hole_id = ?", hole.ID).Find(&floors).Error
	if err != nil {
		return err
	}
	go BulkDelete(Models2IDSlice(floors))

	err = DeleteCache("divisions")
	if err != nil {
		log.Err(err).Msg("DeleteHole: delete cache divisions")
	}

	return c.Status(204).JSON(nil)
}
