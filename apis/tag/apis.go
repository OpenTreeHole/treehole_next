package tag

import (
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ListTags
// @Summary List All Tags
// @Tags Tag
// @Produce application/json
// @Param object query SearchModel false "query"
// @Router /tags [get]
// @Success 200 {array} Tag
func ListTags(c *fiber.Ctx) error {
	var query SearchModel
	err := ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	var tags []Tag
	querySet := DB.Order("temperature DESC")
	if query.Search != "" {
		querySet = querySet.Where("name LIKE ?", "%"+query.Search+"%")
	}
	querySet = querySet.Find(&tags)
	return c.JSON(&tags)
}

// GetTag
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
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// bind and create tag
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
// @Summary Modify A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [put]
// @Param id path int true "id"
// @Param json body ModifyModel true "json"
// @Success 200 {object} Tag
// @Failure 404 {object} MessageModel
func ModifyTag(c *fiber.Ctx) error {
	// validate body
	var body ModifyModel
	err := ValidateBody(c, &body)
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
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}
	MyLog("Tag", "Modify", tag.ID, userID)
	return c.JSON(&tag)
}

// DeleteTag
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
	// validate body
	var body DeleteModel
	err := ValidateBody(c, &body)
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

	err = DB.Transaction(func(tx *gorm.DB) error {
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
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}
	MyLog("Tag", "Delete", id, userID)
	return c.JSON(&newTag)
}
