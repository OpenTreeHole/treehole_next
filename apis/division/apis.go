package division

import (
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AddDivision
// @Summary Add A Division
// @Tags Division
// @Accept application/json
// @Produce application/json
// @Router /divisions [post]
// @Param json body CreateModel true "json"
// @Success 201 {object} models.Division
// @Success 200 {object} models.Division
func AddDivision(c *fiber.Ctx) error {
	var body CreateModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// bind division
	var division Division
	division.Name = body.Name
	division.Description = body.Description
	result := DB.FirstOrCreate(&division, Division{Name: body.Name})
	if result.RowsAffected == 0 {
		c.Status(200)
	} else {
		c.Status(201)
	}
	return Serialize(c, &division)
}

// ListDivisions
// @Summary List All Divisions
// @Tags Division
// @Produce application/json
// @Router /divisions [get]
// @Success 200 {array} models.Division
func ListDivisions(c *fiber.Ctx) error {
	var divisions Divisions
	DB.Find(&divisions)
	return Serialize(c, divisions)
}

// GetDivision
// @Summary Get Division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} models.Division
// @Failure 404 {object} MessageModel
func GetDivision(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var division Division
	result := DB.First(&division, id)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &division)
}

// ModifyDivision
// @Summary Modify A Division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [put]
// @Param id path int true "id"
// @Param json body ModifyModel true "json"
// @Success 200 {object} models.Division
// @Failure 404 {object} MessageModel
func ModifyDivision(c *fiber.Ctx) error {
	var division Division
	var body ModifyModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	id, _ := c.ParamsInt("id")
	division.ID = id
	division.Name = body.Name
	division.Description = body.Description
	division.Pinned = body.Pinned
	if result := DB.Model(&division).Updates(division); result.RowsAffected == 0 { // nothing updated, means that the record does not exist
		return gorm.ErrRecordNotFound
	}
	return Serialize(c, &division)
}

// DeleteDivision
// @Summary Delete A Division
// @Description Delete a division and move all of its holes to another given division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [delete]
// @Param id path int true "id"
// @Param json body DeleteModel true "json"
// @Success 204
// @Failure 404 {object} MessageModel
func DeleteDivision(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var body DeleteModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	if body.To == 0 { // default 1
		body.To = 1
	}
	if id == body.To {
		return BadRequest("The deleted division can't be the same as to.")
	}
	DB.Exec("UPDATE hole SET division_id = ? WHERE division_id = ?", body.To, id)
	DB.Delete(&Division{}, id)
	return c.Status(204).JSON(nil)
}
