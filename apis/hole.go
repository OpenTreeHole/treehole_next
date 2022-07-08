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
// @Failure 500 {object} schemas.MessageModel
func ListHolesByDivision(c *fiber.Ctx) error {
	var query schemas.QueryTime
	err := c.QueryParser(&query)
	if err != nil {
		return err
	}
	if query.Offset.IsZero() {
		query.Offset = time.Now()
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get holes
	var holes Holes
	querySet := Hole{}.MakeQuerySet(query.Offset, query.Size, false)
	if id != 0 {
		querySet = querySet.Where("division_id = ?", id)
	}
	querySet.Find(&holes)

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
	querySet := Hole{}.MakeQuerySet(query.Offset, query.Size, false)
	err = querySet.Model(&tag).
		Association("Holes").Find(&holes)
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
	querySet := Hole{}.MakeQuerySet(query.Offset, query.Size, false)
	if query.Tag != "" {
		var tag Tag
		result := DB.Where("name = ?", query.Tag).First(&tag)
		if result.Error != nil {
			return result.Error
		}
		err = querySet.Model(&tag).Association("Holes").Find(&holes)
		if err != nil {
			return err
		}
	} else if query.DivisionID != 0 {
		querySet.
			Where("division_id = ?", query.DivisionID).
			Find(&holes)
	} else {
		querySet.Find(&holes)
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
	querySet := Hole{}.MakeHiddenQuerySet(false)
	result := querySet.First(&hole, id)
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
	divisionID, _ := c.ParamsInt("id")

	hole := Hole{
		DivisionID: divisionID,
	}
	for _, tag := range body.Tags {
		hole.Tags = append(hole.Tags, &Tag{Name: tag.Name})
	}
	err = hole.Create(c, &body.Content)
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

	hole := Hole{
		DivisionID: body.DivisionID,
	}
	for _, tag := range body.Tags {
		hole.Tags = append(hole.Tags, &Tag{Name: tag.Name})
	}
	err = hole.Create(c, &body.Content)
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
	// validate
	var body schemas.ModifyHole
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	holeID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// Find hole
	var hole Hole
	result := DB.First(&hole, holeID)
	if result.Error != nil {
		return result.Error
	}

	// update divisionID
	if body.DivisionID != 0 {
		hole.DivisionID = body.DivisionID
	}

	// update tags
	if len(body.Tags) != 0 {
		for _, tag := range body.Tags {
			hole.Tags = append(hole.Tags, &Tag{Name: tag.Name})
		}
		err = DB.Transaction(func(tx *gorm.DB) error {
			return hole.SetTags(tx, true)
		})
		if err != nil {
			return err
		}
	}

	// save
	DB.Omit("Tags").Save(&hole)
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
	holeID, _ := c.ParamsInt("id")
	var hole Hole
	hole.ID = holeID
	result := DB.Model(&hole).Select("Hidden").Updates(Hole{Hidden: true})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return c.Status(204).JSON(nil)
}
