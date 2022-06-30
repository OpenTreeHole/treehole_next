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
		err := querySet.Model(&tag).Association("Holes").Find(&holes)
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
	if divisionID == 0 {
		divisionID = 1
	}
	var tagNames []string
	for _, v := range body.Tags {
		tagNames = append(tagNames, v.Name)
	}

	hole, err := createHole(body.Content, body.SpecialTag, tagNames, divisionID)
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
	var tagNames []string
	for _, v := range body.Tags {
		tagNames = append(tagNames, v.Name)
	}
	if body.DivisionID == 0 {
		body.DivisionID = 1
	}

	hole, err := createHole(body.Content, body.SpecialTag, tagNames, body.DivisionID)
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
	holeID, _ := c.ParamsInt("id")

	// Find hole
	var hole Hole
	result := DB.First(&hole, holeID)
	if result.Error != nil {
		return result.Error
	}

	if body.DivisionID != 0 {
		hole.DivisionID = body.DivisionID
	}
	if len(body.Tags) == 0 {
		DB.Save(&hole)
		return Serialize(c, &hole)
	} else {
		var tagNames []string
		for _, v := range body.Tags {
			tagNames = append(tagNames, v.Name)
		}
		err = DB.Transaction(func(tx *gorm.DB) error {
			err := DB.Model(&hole).Association("Tags").Clear()
			if err != nil {
				return err
			}
			err = findCreateTag(&hole, tagNames)
			if err != nil {
				return err
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
	holeID, _ := c.ParamsInt("id")
	var hole Hole
	hole.ID = holeID
	result := DB.Model(&hole).Select("Hidden").Updates(Hole{Hidden: true})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return c.Status(204).JSON(nil)
}

// find or create Tags and increment Tags.temperature according to tagNames
func findCreateTag(hole *Hole, tagNames []string) error {
	result := DB.Where("name IN ?", tagNames).Find(&hole.Tags)
	if result.Error != nil {
		return result.Error
	}
	existTagNames := []string{}
	for _, val := range hole.Tags {
		existTagNames = append(existTagNames, val.Name)
	}
	createTagNames := DiffrenceSet(tagNames, existTagNames)
	for _, val := range createTagNames {
		hole.Tags = append(hole.Tags, &Tag{Name: val})
	}
	for i := 0; i < len(hole.Tags); i++ {
		hole.Tags[i].Temperature += 1
	}
	return nil
}

func createHole(content string, specialTag string,
	tagNames []string, divisionID int) (hole Hole, err error) {

	// Bind and Create floor
	var floor Floor
	floor.Content = content
	floor.SpecialTag = specialTag
	hole.DivisionID = divisionID
	hole.Floors = []Floor{floor}

	// Find and create Tags
	err = findCreateTag(&hole, tagNames)
	if err != nil {
		return Hole{}, err
	}

	// Create hole
	result := DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&hole)
	if result.Error != nil {
		return Hole{}, result.Error
	}

	return hole, nil
}
