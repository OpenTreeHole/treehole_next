package division

import (
	"github.com/gofiber/fiber/v2"
	"treehole_next/db"
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
	result := db.DB.Where(&Division{Name: division.Name}).FirstOrCreate(&division)
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
	db.DB.Find(&divisions)
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
	if result := db.DB.First(&division, id); result.Error != nil {
		return result.Error
	}
	return c.JSON(division)
}
