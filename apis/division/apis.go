package division

import (
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"

	"github.com/opentreehole/go-common"

	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

// AddDivision
//
// @Summary Add A Division
// @Tags Division
// @Accept application/json
// @Produce application/json
// @Router /divisions [post]
// @Param json body CreateModel true "json"
// @Success 201 {object} models.Division
// @Success 200 {object} models.Division
func AddDivision(c *fiber.Ctx) error {
	// validate body
	var body CreateModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission check
	if !user.IsAdmin {
		return common.Forbidden()
	}

	// bind division
	division := Division{
		Name:        body.Name,
		Description: body.Description,
	}
	result := DB.FirstOrCreate(&division, Division{Name: body.Name})
	if result.RowsAffected == 0 {
		c.Status(200)
	} else {
		c.Status(201)
	}
	return Serialize(c, &division)
}

// ListDivisions
//
// @Summary List All Divisions
// @Tags Division
// @Produce application/json
// @Router /divisions [get]
// @Success 200 {array} models.Division
func ListDivisions(c *fiber.Ctx) error {
	var divisions Divisions
	if GetCache("divisions", &divisions) {
		return c.JSON(divisions)
	}
	err := DB.Find(&divisions, "hidden = false").Error
	if err != nil {
		return err
	}
	return Serialize(c, divisions)
}

// GetDivision
//
// @Summary Get Division
// @Tags Division
// @Produce application/json
// @Router /divisions/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} models.Division
// @Failure 404 {object} MessageModel
func GetDivision(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var division Division
	result := DB.Where("hidden = false").First(&division, id)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &division)
}

// ModifyDivision
//
// @Summary Modify A Division
// @Tags Division
// @Produce json
// @Router /divisions/{id} [put]
// @Router /divisions/{id}/_modify [patch]
// @Param id path int true "id"
// @Param json body ModifyDivisionModel true "json"
// @Success 200 {object} models.Division
// @Failure 404 {object} MessageModel
func ModifyDivision(c *fiber.Ctx) error {
	// validate body
	var body ModifyDivisionModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission check
	if !user.IsAdmin {
		return common.Forbidden()
	}

	var division Division
	err = DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&division, id).Error
		if err != nil {
			return err
		}

		modifyData := make(map[string]any)
		if body.Name != nil {
			modifyData["name"] = *body.Name
		}
		if body.Description != nil {
			modifyData["description"] = *body.Description
		}
		if body.Pinned != nil {
			data, _ := json.Marshal(body.Pinned)
			modifyData["pinned"] = string(data)
		}

		if len(modifyData) == 0 {
			return common.BadRequest("No data to modify.")
		}

		return tx.Model(&division).Updates(modifyData).Error
	})
	if err != nil {
		return err
	}

	var newDivision Division
	err = DB.First(&newDivision, id).Error
	if err != nil {
		return err
	}

	MyLog("Division", "Modify", division.ID, user.ID, RoleAdmin)

	CreateAdminLog(DB, AdminLogTypeDivision, user.ID, map[string]any{
		"division_id": division.ID,
		"before":      division,
		"after":       newDivision,
	})

	// refresh cache. here should not use `go refreshCache`
	err = refreshCache(c)
	if err != nil {
		return err
	}

	return Serialize(c, &newDivision)
}

// DeleteDivision
//
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
	// validate body
	var body DeleteModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return common.Forbidden()
	}

	if id == body.To {
		return common.BadRequest("The deleted division can't be the same as to.")
	}
	err = DB.Exec("UPDATE hole SET division_id = ? WHERE division_id = ?", body.To, id).Error
	if err != nil {
		return err
	}
	err = DB.Delete(&Division{ID: id}).Error
	if err != nil {
		return err
	}

	// log
	//if err != nil {
	//	return err
	//}
	MyLog("Division", "Delete", id, user.ID, RoleAdmin, "To: ", strconv.Itoa(body.To))

	err = refreshCache(c)
	if err != nil {
		return err
	}

	return c.Status(204).JSON(nil)
}
