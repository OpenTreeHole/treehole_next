package tag

import (
	"strings"

	"github.com/opentreehole/go-common"
	"gorm.io/plugin/dbresolver"

	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ListTags
//
// @Summary List All Tags
// @Tags Tag
// @Produce application/json
// @Param object query SearchModel false "query"
// @Router /tags [get]
// @Success 200 {array} Tag
func ListTags(c *fiber.Ctx) error {
	var query SearchModel
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	tags := make(Tags, 0, 10)
	if query.Search == "" {
		if GetCache("tags", &tags) {
			return c.JSON(&tags)
		} else {
			err = DB.Order("temperature DESC").Find(&tags).Error
			if err != nil {
				return err
			}
			go UpdateTagCache(tags)
			return c.JSON(&tags)
		}
	}
	err = DB.Where("name LIKE ?", "%"+query.Search+"%").
		Order("temperature DESC").Find(&tags).Error
	if err != nil {
		return err
	}
	return c.JSON(&tags)
}

// GetTag
//
// @Summary Get A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Tag
// @Failure 404 {object} MessageModel
func GetTag(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var tag Tag
	tag.ID = id
	result := DB.First(&tag)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&tag)
}

// CreateTag
//
// @Summary Create A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags [post]
// @Param json body CreateModel true "json"
// @Success 200 {object} Tag
// @Success 201 {object} Tag
func CreateTag(c *fiber.Ctx) error {
	// validate body
	var tag Tag
	var body CreateModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// check tag prefix
	user, err := GetUser(c)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		if strings.HasPrefix(body.Name, "#") {
			return common.BadRequest("只有管理员才能创建 # 开头的 tag")
		}
		if strings.HasPrefix(body.Name, "@") {
			return common.BadRequest("只有管理员才能创建 @ 开头的 tag")
		}
	}

	// bind and create tag
	body.Name = strings.TrimSpace(body.Name)
	tag.Name = body.Name
	result := DB.Where("name = ?", body.Name).FirstOrCreate(&tag)
	if result.RowsAffected == 0 {
		c.Status(200)
	} else {
		c.Status(201)
	}
	return c.JSON(&tag)
}

// ModifyTag
//
// @Summary Modify A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [put]
// @Param id path int true "id"
// @Param json body ModifyModel true "json"
// @Success 200 {object} Tag
// @Failure 404 {object} MessageModel
func ModifyTag(c *fiber.Ctx) error {
	// admin
	user, err := GetUser(c)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return common.Forbidden()
	}

	// validate body
	var body ModifyModel
	err = common.ValidateBody(c, &body)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// modify tag
	var tag Tag
	DB.Find(&tag, id)
	tag.Name = body.Name
	tag.Temperature = body.Temperature
	DB.Save(&tag)

	// log
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}
	MyLog("Tag", "Modify", tag.ID, userID, RoleAdmin)
	return c.JSON(&tag)
}

// DeleteTag
//
// @Summary Delete A Tag
// @Description Delete a tag and link all of its holes to another given tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [delete]
// @Param id path int true "id"
// @Param json body DeleteModel true "json"
// @Success 200 {object} Tag
// @Failure 404 {object} MessageModel
func DeleteTag(c *fiber.Ctx) error {
	// admin
	user, err := GetUser(c)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return common.Forbidden()
	}

	// validate body
	var body DeleteModel
	err = common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var tag Tag
	result := DB.First(&tag, id)
	if result.Error != nil {
		return result.Error
	}

	var newTag Tag
	result = DB.Where("name = ?", body.To).First(&newTag)
	if result.Error != nil {
		return result.Error
	}

	newTag.Temperature += tag.Temperature

	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		result = tx.Exec(`
 DELETE FROM hole_tags WHERE tag_id = ? AND hole_id IN
 (SELECT a.hole_id FROM
 (SELECT hole_id FROM hole_tags WHERE tag_id = ?)a
 )`, id, newTag.ID)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Exec(`UPDATE hole_tags SET tag_id = ? WHERE tag_id = ?`, newTag.ID, id)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Updates(&newTag)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Delete(&tag)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		return err
	}

	// log
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}
	MyLog("Tag", "Delete", id, userID, RoleAdmin)
	return c.JSON(&newTag)
}
