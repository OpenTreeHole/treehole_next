package apis

import (
	"time"
	. "treehole_next/models"
	"treehole_next/schemas"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ListHolesByDivision
// @Summary List Holes In A Division
// @Tags Hole
// @Produce application/json
// @Router /divisions/{division_id}/holes [get]
// @Param division_id path int true "division_id"
// @Param object query schemas.QueryTime false "query"
// @Success 200 {array} Hole
// @Failure 404 {object} schemas.MessageModel
func ListHolesByDivision(c *fiber.Ctx) error {
	var query schemas.QueryTime
	err := c.QueryParser(&query)
	if err != nil {
		return err
	}
	if query.Offset.IsZero() {
		query.Offset = time.Now()
	}
	id, _ := c.ParamsInt("id")

	// get division
	var division Division
	result := DB.First(&division, id)
	if result.Error != nil {
		return result.Error
	}

	// get holes
	var holes Holes
	result = DB.
		Where("division_id = ?", id).
		Where("updated_at < ?", query.Offset).
		Order("updated_at desc").Limit(query.Size).
		Find(&holes)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return Serialize(c, &holes)
}

// ListHolesByTag
// @Summary List Holes By Tag
// @Tags Hole
// @Produce application/json
// @Router /tags/{tag_name}/holes [get]
// @Param tag_name path string true "tag_name"
// @Param object query schemas.QueryTime false "query"
// @Success 200 {array} Hole
// @Failure 404 {object} schemas.MessageModel
func ListHolesByTag(c *fiber.Ctx) error {
	var query schemas.QueryTime
	err := c.QueryParser(&query)
	if err != nil {
		return err
	}
	if query.Offset.IsZero() {
		query.Offset = time.Now()
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
	err = DB.Model(&tag).
		Where("updated_at < ?", query.Offset).
		Order("updated_at desc").Limit(query.Size).
		Association("Holes").
		Find(&holes)
	if err != nil {
		return err
	}

	return Serialize(c, &holes)
}

// ListHolesOld
// @Summary Old API for Listing Holes
// @Deprecated
// @Tags Hole
// @Produce application/json
// @Router /holes [get]
// @Param object query schemas.GetHoleOld false "query"
// @Success 200 {array} Hole
func ListHolesOld(c *fiber.Ctx) error {
	var query schemas.GetHoleOld
	err := c.QueryParser(&query)
	if err != nil {
		return err
	}
	if query.Offset.IsZero() {
		query.Offset = time.Now()
	}

	var holes Holes
	if query.Tag != "" {
		var tag Tag
		result := DB.Where("name = ?", query.Tag).First(&tag)
		if result.Error != nil {
			return result.Error
		}
		err := DB.Model(&tag).Association("Holes").Find(&holes)
		if err != nil {
			return err
		}
	} else {
		result := DB.
			Where("updated_at < ?", query.Offset).
			Order("updated_at desc").Limit(query.Size).
			Where("division_id = ?", query.DivisionID).
			Find(&holes)
		if result.Error != nil {
			return result.Error
		}
	}

	return Serialize(c, &holes)
}

// GetHole
// @Summary Get A Hole
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Hole
// @Failure 404 {object} schemas.MessageModel
func GetHole(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")

	// get hole
	var hole Hole
	result := DB.First(&hole, id)
	if result.Error != nil {
		return result.Error
	}

	return Serialize(c, &hole)
}

// CreateHole
// @Summary Create A Hole
// @Description Create a hole, create tags and floor binding to it and set the name mapping
// @Tags Hole
// @Produce application/json
// @Router /divisions/{division_id}/holes [post]
// @Param division_id path int true "division id"
// @Param json body schemas.CreateHole true "json"
// @Success 201 {object} Hole
func CreateHole(c *fiber.Ctx) error {
	var body schemas.CreateHole
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	division_id, _ := c.ParamsInt("id")

	// bind hole
	var floor Floor
	floor.Content = body.Content
	floor.SpecialTag = body.SpecialTag
	hole := Hole{
		Tags: make([]*Tag, len(body.Tags)),
	}
	hole.DivisionID = division_id
	hole.Floors = []Floor{floor}

	// Create
	err = DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(hole.Tags); i++ {
			result := tx.Where(Tag{Name: body.Tags[i].Name}).
				FirstOrCreate(&hole.Tags[i])
			if result.Error != nil {
				return result.Error
			}
			hole.Tags[i].Temperature += 1
		}
		result := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&hole)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}

	return Serialize(c.Status(201), &hole)
}

// CreateHoleOld
// @Summary Old API for Creating A Hole
// @Deprecated
// @Tags Hole
// @Produce application/json
// @Router /holes [post]
// @Param json body schemas.CreateHoleOld true "json"
// @Success 201 {object} Hole
func CreateHoleOld(c *fiber.Ctx) error {
	var body schemas.CreateHoleOld
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	var floor Floor
	floor.Content = body.Content
	floor.SpecialTag = body.SpecialTag
	hole := Hole{
		Tags: make([]*Tag, len(body.Tags)),
	}
	hole.DivisionID = body.DivisionID
	hole.Floors = []Floor{floor}

	// Create
	err = DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(hole.Tags); i++ {
			result := tx.Where(Tag{Name: body.Tags[i].Name}).
				FirstOrCreate(&hole.Tags[i])
			if result.Error != nil {
				return result.Error
			}
			hole.Tags[i].Temperature += 1
		}
		result := tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(&hole)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}

	return Serialize(c.Status(201), &hole)
}

// ModifyHole
// @Summary Modify A Hole
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [put]
// @Param id path int true "id"
// @Param json body schemas.ModifyHole true "json"
// @Success 200 {object} Hole
// @Failure 404 {object} schemas.MessageModel
func ModifyHole(c *fiber.Ctx) error {
	var body schemas.ModifyHole
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}
	hole_id, _ := c.ParamsInt("id")
	var hole Hole
	result := DB.Where("id = ?", hole_id).First(&hole)
	if result.Error != nil {
		return result.Error
	}
	hole.DivisionID = body.DivisionID

	// Modify
	err = DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(hole.Tags); i++ {
			result := tx.Where(Tag{Name: body.Tags[i].Name}).
				FirstOrCreate(&hole.Tags[i])
			if result.Error != nil {
				return result.Error
			}
			hole.Tags[i].Temperature += 1
		}
		result := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&hole)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}

	return Serialize(c, &hole)
}

// DeleteHole
// @Summary Delete A Hole
// @Description Hide a hole, but visible to admins. This may affect many floors, DO NOT ABUSE!!!
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [delete]
// @Param id path int true "id"
// @Success 204
// @Failure 404 {object} schemas.MessageModel
func DeleteHole(c *fiber.Ctx) error {
	hole_id, _ := c.ParamsInt("id")
	var hole Hole
	result := DB.Where("id = ?", hole_id).First(&hole)
	if result.Error != nil {
		c.Status(404)
		return result.Error
	}
	hole.Hidden = true
	// Modify
	err := DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&hole)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		c.Status(404)
		return err
	}
	return c.Status(204).JSON(nil)
}
