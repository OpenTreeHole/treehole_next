package hole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"
	"treehole_next/config"
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

// ListHoles
//
// @Summary API for Listing Holes
// @Tags Hole
// @Produce json
// @Router /holes [get]
// @Param object query ListOldModel false "query"
// @Success 200 {array} Hole
func ListHoles(c *fiber.Ctx) error {
	var query ListOldModel
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	var holes Holes
	err = DB.Transaction(func(tx *gorm.DB) error {
		querySet, err := holes.MakeQuerySet(query.Offset, query.Size, query.Order, c)
		if err != nil {
			return err
		}

		if query.CreatedStart != nil {
			querySet = querySet.Where("hole.created_at >= ?", query.CreatedStart.Time)
		}

		if query.CreatedEnd != nil {
			querySet = querySet.Where("hole.created_at <= ?", query.CreatedEnd.Time)
		}

		if query.DivisionID != 0 {
			querySet = querySet.Where("hole.division_id = ?", query.DivisionID)
		}

		if len(query.Tags) != 0 {
			var tags []Tag
			err = DB.Where("name IN ?", query.Tags).Find(&tags).Error
			if err != nil {
				return err
			}

			if len(tags) != len(query.Tags) {
				return common.BadRequest("部分标签不存在")
			}

			tagIDs := make([]int, len(tags))
			for i, tag := range tags {
				tagIDs[i] = tag.ID
			}

			var holeIDs []int
			err = DB.Table("hole_tags").
				Select("hole_id").
				Where("tag_id IN ?", tagIDs).
				Group("hole_id").
				Having("COUNT(DISTINCT tag_id) = ?", len(tagIDs)).
				Pluck("hole_id", &holeIDs).Error
			if err != nil {
				return err
			}

			querySet = querySet.Where("hole.id IN ?", holeIDs)
			err = querySet.Find(&holes).Error
			if err != nil {
				return err
			}
		} else if query.Tag != "" {
			var tag Tag
			err = DB.Where("name = ?", query.Tag).Find(&tag).Error
			if err != nil {
				return err
			}
			err = querySet.Model(&tag).Order("updated_at desc").Association("Holes").Find(&holes)
			if err != nil {
				return err
			}
		} else {
			err = querySet.Find(&holes).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
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
	user, err := GetCurrLoginUser(c)
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
	user, err := GetCurrLoginUser(c)
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
// @Router /holes/{id}/_webvpn [patch]
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
	user, err := GetCurrLoginUser(c)
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
		if body.Hidden != nil {
			if *body.Hidden && !hole.Hidden {
				hole.Hidden = true
				changed = true

				// delete floors from Elasticsearch
				var floors Floors
				_ = DB.Where("hole_id = ?", hole.ID).Find(&floors)
				go BulkDelete(Models2IDSlice(floors))

				// log
				MyLog("Hole", "Modify", holeID, user.ID, RoleAdmin, "Hidden: ")
			} else if !*body.Hidden && hole.Hidden {
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
		} else {
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

		// modify frozen
		if body.Frozen != nil {
			changed = true
			hole.Frozen = *body.Frozen

			MyLog("Hole", "Modify", holeID, user.ID, RoleAdmin, "Frozen: ")
		}

		// save
		if changed {
			err = tx.Model(&hole).
				Omit(clause.Associations, "UpdatedAt").
				Select("DivisionID", "Hidden", "Locked", "Frozen").
				Updates(&hole).Error
			if err != nil {
				return err
			}

			if user.IsAdmin {
				CreateAdminLog(tx, AdminLogTypeHole, user.ID, struct {
					HoleID int            `json:"hole_id"`
					Before map[string]any `json:"before"`
					Modify ModifyModel    `json:"modify"`
				}{
					HoleID: holeID,
					Before: map[string]any{
						"division_id": hole.DivisionID,
						"hidden":      hole.Hidden,
						"locked":      hole.Locked,
						"frozen":      hole.Frozen,
						"tags":        hole.Tags,
					},
					Modify: body,
				})
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
	user, err := GetCurrLoginUser(c)
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
	CreateAdminLog(DB, AdminLogTypeHideHole, user.ID, struct {
		HoleID int  `json:"hole_id"`
		Hidden bool `json:"hidden"`
	}{
		HoleID: holeID,
		Hidden: true,
	})

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

	user, err := GetCurrLoginUser(c)
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

	result := DB.Delete(&hole)
	if result.Error != nil {
		return result.Error
	}
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
func GenerateSummary(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")
	var cachedData Summary
	if GetCache("AISummary"+strconv.Itoa(id), &cachedData) {
		switch cachedData.Code {
		case 1000:
			return c.Status(200).JSON(cachedData)
		case 1001:
			// get new summary
			resp, err := http.Get(config.Config.AISummaryURL + "/get_summary?code=1000&hole_id=" + strconv.Itoa(id))
			if err != nil {
				log.Err(err).Msg("AISummary: get summary from server err")
				return err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			err = json.Unmarshal(body, &cachedData)
			if err != nil {
				return err
			}
			if len(cachedData.Data.Interactions) > 5 {
				cachedData.Data.Interactions = cachedData.Data.Interactions[:5]
			}

			// renew cache
			if cachedData.Code == 1000 || cachedData.Code == 1001 {
				//sensitiveCheckResp, err := sensitive.CheckSensitive(sensitive.ParamsForCheck{
				//	Content:  cachedData.Data.Summary,
				//	Id:       time.Now().UnixNano(),
				//	TypeName: sensitive.TypeFloor,
				//})
				if true {
					// set default interaction_type
					for i := range cachedData.Data.Interactions {
						if cachedData.Data.Interactions[i].InteractionType == "" {
							cachedData.Data.Interactions[i].InteractionType = "reply"
						}
					}
					err = SetCache("AISummary"+strconv.Itoa(id), cachedData, 24*time.Hour)
					if err != nil {
						log.Err(err).Msg("AISummary: set cache err")
					}
				}

			} else {
				err := DeleteCache("AISummary" + strconv.Itoa(id))
				if err != nil {
					log.Err(err).Msg("AISummary: delete cache err")
				}
			}

			return c.Status(200).JSON(cachedData)
		default:
			err := DeleteCache("AISummary" + strconv.Itoa(id))
			if err != nil {
				log.Err(err).Msg("AISummary: delete cache err")
			}
		}
	}
	// if no cache or the data of cache is invalid, generate a new summary

	// get hole
	holeSet, err := MakeHoleQuerySet(c)
	if err != nil {
		return err
	}
	var hole Hole
	err = holeSet.Take(&hole, id).Error
	if err != nil {
		return err
	}
	err = hole.Preprocess(c)
	if err != nil {
		return err
	}

	if !hole.AISummaryAvailable {
		return c.Status(200).JSON(fiber.Map{
			"code":    2002,
			"message": "unavailable",
			"data":    fiber.Map{},
		})
	}

	// get floors
	var floors Floors
	floorSet, err := floors.MakeQuerySet(&id, nil, nil, c)
	if err != nil {
		return err
	}

	err = floorSet.Find(&floors).Error
	if err != nil {
		return err
	}

	content := ""
	if hole.HoleFloor.FirstFloor != nil {
		content = hole.HoleFloor.FirstFloor.Content
	}

	requestBody := map[string]any{
		"floors":  floors,
		"content": content,
		"hole_id": hole.ID,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	var response struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Data    struct{} `json:"data"`
	}

	resp, err := http.Post(config.Config.AISummaryURL+"/generate_summary?code=1000", "application/json", bytes.NewReader(requestJSON))
	if err != nil {
		log.Err(err).Msg("AISummary: generate summary from server err")
		response.Code = 3001
		response.Message = "service_error"
		return c.Status(200).JSON(response)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	//create a cache when generate a new summary
	switch response.Code {
	case 1000, 1001:
		var cache Summary
		cache.Code = 1001
		err := SetCache("AISummary"+strconv.Itoa(id), cache, 24*time.Hour)
		if err != nil {
			log.Err(err).Msg("AISummary: set cache err")
		}

		response.Code = 1002
		response.Message = "started"
	case 1002, 2001, 2002, 3001, 3002:

	default:
		response.Code = 3001
		response.Message = "service_error"
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	return c.Status(200).JSON(response)
}

func GetFeedback(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{})
}
