package apis

import (
	"github.com/gofiber/fiber/v2"
	. "treehole_next/models"
)

// AddDivision
// @Summary Add A Division
// @Tags Division
// @Accept application/json
// @Produce application/json
// @Router /divisions [post]
// @Param json body AddDivisionModel true "json"
// @Success 201 {object} Division
// @Success 200 {object} Division
func AddDivision(c *fiber.Ctx) error {
	var division Division
	if err := c.BodyParser(&division); err != nil {
		return err
	}
	result := DB.Where("name = ?", division.Name).FirstOrCreate(&division)
	if result.RowsAffected == 0 {
		c.Status(200)
	} else {
		c.Status(201)
	}
	return c.JSON(division)
}

// ListDivisions
// @Summary List Divisions
// @Tags Division
// @Produce application/json
// @Router /divisions [get]
// @Success 200 {array} Division
func ListDivisions(c *fiber.Ctx) error {
	var divisions []Division
	DB.Find(&divisions)
	return c.JSON(divisions)
}

// GetDivision
// @Summary Get Division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Division
// @Failure 404 {object} utils.MessageModel
func GetDivision(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var division Division
	if result := DB.First(&division, id); result.Error != nil {
		return result.Error
	}
	return c.JSON(division)
}

// ModifyDivision
// @Summary Modify A Division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [put]
// @Param id path int true "id"
// @Param json body ModifyDivisionModel true "json"
// @Success 200 {object} Division
// @Failure 404 {object} utils.MessageModel
func ModifyDivision(c *fiber.Ctx) error {
	return nil
}

// DeleteDivision
// @Summary Delete A Division
// @Description Delete a division and move all of its holes to another given division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [delete]
// @Param id path int true "id"
// @Param json body DeleteDivisionModel true "json"
// @Success 200 {object} Division
// @Failure 404 {object} utils.MessageModel
func DeleteDivision(c *fiber.Ctx) error {
	return nil
}
