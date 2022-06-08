package apis

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	. "treehole_next/models"
	"treehole_next/schemas"
	. "treehole_next/utils"
)

func serializeDivision(c *fiber.Ctx, division *Division) error {
	var divisionResponse schemas.DivisionResponse

	// save division.Pinned and remove it to ensure conversion
	var pinned []int = division.Pinned
	division.Pinned = nil
	// merge division into divisionResponse
	data, err := json.Marshal(&division)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &divisionResponse)
	if err != nil {
		return err
	}

	// set divisionResponse.Pinned
	var holes []Hole
	DB.Find(&holes, pinned)
	orderedHoles := make([]Hole, 0, len(holes))
	for _, order := range pinned {
		// binary search the index
		index := func(target int) int {
			left := 0
			right := len(holes)
			for left < right {
				mid := left + (right-left)>>1
				if holes[mid].ID < target {
					left = mid + 1
				} else if holes[mid].ID > target {
					right = mid
				} else {
					return mid
				}
			}
			return -1
		}(order)
		if index >= 0 {
			orderedHoles = append(orderedHoles, holes[index])
		}
	}
	divisionResponse.Pinned = orderedHoles

	return c.JSON(divisionResponse)
}

// AddDivision
// @Summary Add A Division
// @Tags Division
// @Accept application/json
// @Produce application/json
// @Router /divisions [post]
// @Param json body schemas.AddDivisionModel true "json"
// @Success 201 {object} schemas.DivisionResponse
// @Success 200 {object} schemas.DivisionResponse
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
	return serializeDivision(c, &division)
}

// ListDivisions
// @Summary List All Divisions
// @Tags Division
// @Produce application/json
// @Router /divisions [get]
// @Success 200 {array} schemas.DivisionResponse
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
// @Success 200 {object} schemas.DivisionResponse
// @Failure 404 {object} schemas.MessageModel
func GetDivision(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var division Division
	if result := DB.First(&division, id); result.Error != nil {
		return result.Error
	}
	return serializeDivision(c, &division)
}

// ModifyDivision
// @Summary Modify A Division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [put]
// @Param id path int true "id"
// @Param json body schemas.ModifyDivisionModel true "json"
// @Success 200 {object} schemas.DivisionResponse
// @Failure 404 {object} schemas.MessageModel
func ModifyDivision(c *fiber.Ctx) error {
	var division Division
	if err := c.BodyParser(&division); err != nil {
		return err
	}
	id, _ := c.ParamsInt("id")
	division.ID = id
	result := DB.Model(&division).Updates(division)
	if result.RowsAffected == 0 { // nothing updated, means that the record does not exist
		return gorm.ErrRecordNotFound
	}
	return serializeDivision(c, &division)
}

// DeleteDivision
// @Summary Delete A Division
// @Description Delete a division and move all of its holes to another given division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [delete]
// @Param id path int true "id"
// @Param json body schemas.DeleteDivisionModel true "json"
// @Success 204
// @Failure 404 {object} schemas.MessageModel
func DeleteDivision(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var body schemas.DeleteDivisionModel
	if err := BindJSON(c, &body); err != nil {
		return err
	}
	if body.To == 0 { // default 1
		body.To = 1
	}
	if id == body.To {
		return BadRequest("The deleted division can't be the same as to.")
	}
	DB.Exec("update hole set division_id = ? where division_id = ?", body.To, id)
	DB.Delete(&Division{}, id)
	return c.Status(204).JSON(nil)
}
