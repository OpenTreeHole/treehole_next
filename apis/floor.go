package apis

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	. "treehole_next/models"
)

// ListFloorsInAHole
// @Summary List Floors In A Hole
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [get]
// @Param hole_id path int true "hole id"
// @Param object query schemas.Query false "query"
// @Success 200 {array} Floor
func ListFloorsInAHole(c *fiber.Ctx) error {
	fmt.Println(Floor{})
	return nil
}

// ListFloorsOld
// @Summary Old API for Listing Floors
// @Deprecated
// @Tags Floor
// @Produce application/json
// @Router /floors [get]
// @Param object query schemas.ListFloorOld false "query"
// @Success 200 {array} Floor
func ListFloorsOld(c *fiber.Ctx) error {
	return nil
}

// GetFloor
// @Summary Get A Floor
// @Tags Floor
// @Produce application/json
// @Router /floors/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Floor
// @Failure 404 {object} schemas.MessageModel
func GetFloor(c *fiber.Ctx) error {
	return nil
}

// CreateFloor
// @Summary Create A Floor
// @Tags Floor
// @Produce application/json
// @Router /holes/{hole_id}/floors [post]
// @Param hole_id path int true "hole id"
// @Param json body schemas.CreateFloor true "json"
// @Success 201 {object} Floor
func CreateFloor(c *fiber.Ctx) error {
	return nil
}

// CreateFloorOld
// @Summary Old API for Creating A Floor
// @Deprecated
// @Tags Floor
// @Produce application/json
// @Router /floors [post]
// @Param json body schemas.CreateFloorOld true "json"
// @Success 201 {object} Floor
func CreateFloorOld(c *fiber.Ctx) error {
	return nil
}

// ModifyFloor
// @Summary Modify A Floor
// @Tags Floor
// @Produce application/json
// @Router /floors/{id} [put]
// @Param id path int true "id"
// @Param json body schemas.ModifyFloor true "json"
// @Success 200 {object} Floor
// @Failure 404 {object} schemas.MessageModel
func ModifyFloor(c *fiber.Ctx) error {
	return nil
}

// DeleteFloor
// @Summary Delete A Floor
// @Tags Floor
// @Produce application/json
// @Router /floors/{id} [delete]
// @Param id path int true "id"
// @Success 200 {object} Floor
// @Failure 404 {object} schemas.MessageModel
func DeleteFloor(c *fiber.Ctx) error {
	return nil
}
