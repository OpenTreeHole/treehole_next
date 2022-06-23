package apis

import (
	. "treehole_next/models"
	"treehole_next/schemas"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

// ListTags
// @Summary List All Tags
// @Tags Tag
// @Produce application/json
// @Router /tags [get]
// @Success 200 {array} Tag
func ListTags(c *fiber.Ctx) error {
	var tags []*Tag
	DB.Find(&tags)
	return c.JSON(&tags)
}

// GetTag
// @Summary Get A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Tag
// @Failure 404 {object} schemas.MessageModel
func GetTag(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var tag Tag
	tag.ID = id
	if result := DB.First(&tag); result.Error != nil {
		return result.Error
	}
	return c.JSON(&tag)
}

// CreateTag
// @Summary Create A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags [post]
// @Param json body schemas.CreateTag true "json"
// @Success 200 {object} Tag
// @Success 201 {object} Tag
func CreateTag(c *fiber.Ctx) error {
	var tag Tag
	var body schemas.CreateTag
	if err := c.BodyParser(&body); err != nil {
		return err
	}
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
// @Param json body schemas.ModifyTag true "json"
// @Success 200 {object} Tag
// @Failure 404 {object} schemas.MessageModel
func ModifyTag(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var tag Tag
	var body schemas.ModifyTag
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	DB.Find(&tag, id)
	tag.Name = body.Name
	tag.Temperature = body.Temperature
	DB.Save(&tag)
	return c.JSON(&tag)
}

// DeleteTag
// @Summary Delete A Tag
// @Description Delete a tag and link all of its holes to another given tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [delete]
// @Param id path int true "id"
// @Param json body schemas.DeleteTag true "json"
// @Success 200 {object} Tag
// @Failure 404 {object} schemas.MessageModel
func DeleteTag(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var tag Tag
	var newtag Tag
	var body schemas.DeleteTag

	if result := DB.First(&tag, id); result.Error != nil {
		return result.Error
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	newtag.Name = body.To
	if result := DB.Where("name = ?", newtag.Name).First(&newtag); result.Error != nil {
		return result.Error
	}

	newtag.Temperature += tag.Temperature

	err := DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
			DELETE FROM hole_tags WHERE tag_id = ? AND hole_id IN
				(SELECT a.hole_id FROM
					(SELECT hole_id FROM hole_tags WHERE tag_id = ?)a
			)`, id, newtag.ID)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Exec(`UPDATE hole_tags SET tag_id = ? WHERE tag_id = ?`, newtag.ID, id)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Updates(&newtag)
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
	return c.JSON(&newtag)
}
