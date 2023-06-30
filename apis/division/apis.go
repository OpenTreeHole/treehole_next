package division

import (
	"github.com/opentreehole/go-common"
	"strconv"
	"treehole_next/config"
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AddDivision
//
//	@Summary	Add A Division
//	@Tags		Division
//	@Accept		application/json
//	@Produce	application/json
//	@Router		/divisions [post]
//	@Param		json	body		CreateModel	true	"json"
//	@Success	201		{object}	models.Division
//	@Success	200		{object}	models.Division
func AddDivision(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[CreateModel](c)
	if err != nil {
		return err
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
//	@Summary	List All Divisions
//	@Tags		Division
//	@Produce	application/json
//	@Router		/divisions [get]
//	@Success	200	{array}	models.Division
func ListDivisions(c *fiber.Ctx) error {
	var divisions Divisions
	if GetCache("divisions", &divisions) {
		return c.JSON(divisions)
	}
	DB.Find(&divisions)
	return Serialize(c, divisions)
}

// GetDivision
//
//	@Summary	Get Division
//	@Tags		Division
//	@Produce	application/json
//	@Router		/divisions/{id} [get]
//	@Param		id	path		int	true	"id"
//	@Success	200	{object}	models.Division
//	@Failure	404	{object}	MessageModel
func GetDivision(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var division Division
	result := DB.First(&division, id)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &division)
}

// ModifyDivision
//
//	@Summary	Modify A Division
//	@Tags		Division
//	@Produce	json
//	@Router		/divisions/{id} [put]
//	@Param		id		path		int			true	"id"
//	@Param		json	body		ModifyModel	true	"json"
//	@Success	200		{object}	models.Division
//	@Failure	404		{object}	MessageModel
func ModifyDivision(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[ModifyModel](c)
	if err != nil {
		return err
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	division := Division{
		Name:        body.Name,
		Description: body.Description,
		Pinned:      body.Pinned,
	}
	division.ID = id
	result := DB.Model(&division).Updates(division)
	// nothing updated, means that the record does not exist
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	// log
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	MyLog("Division", "Modify", division.ID, userID, RoleAdmin)

	if config.Config.Mode != "test" {
		go refreshCache()
	} else {
		refreshCache()
	}

	return Serialize(c, &division)
}

// DeleteDivision
//
//	@Summary		Delete A Division
//	@Description	Delete a division and move all of its holes to another given division
//	@Tags			Division
//	@Produce		application/json
//	@Router			/divisions/{id} [delete]
//	@Param			id		path	int			true	"id"
//	@Param			json	body	DeleteModel	true	"json"
//	@Success		204
//	@Failure		404	{object}	MessageModel
func DeleteDivision(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[DeleteModel](c)
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
	if err != nil {
		return err
	}
	MyLog("Division", "Delete", id, user.ID, RoleAdmin, "To: ", strconv.Itoa(body.To))

	go refreshCache()

	return c.Status(204).JSON(nil)
}
