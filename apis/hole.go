package apis

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	. "treehole_next/models"
)

// ListHolesByDivision
// @Summary List Holes In A Division
// @Tags Hole
// @Produce application/json
// @Router /divisions/{division_id}/holes [get]
// @Param division_id path int true "division_id"
// @Param object query schemas.DocQuery false "query"
// @Success 200 {array} Hole
func ListHolesByDivision(c *fiber.Ctx) error {
	fmt.Println(Hole{})
	return nil
}

// ListHolesByTag
// @Summary List Holes By Tag
// @Tags Hole
// @Produce application/json
// @Router /tags/{tag_name}/holes [get]
// @Param tag_name path string true "tag_name"
// @Param object query schemas.DocQuery false "query"
// @Success 200 {array} Hole
func ListHolesByTag(c *fiber.Ctx) error {
	fmt.Println(Hole{})
	return nil
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
	fmt.Println(Hole{})
	return nil
}

// GetHole
// @Summary Get A Hole
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Hole
// @Failure 404 {object} utils.MessageModel
func GetHole(c *fiber.Ctx) error {
	return nil
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
	return nil
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
	return nil
}

// ModifyHole
// @Summary Modify A Hole
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [put]
// @Param id path int true "id"
// @Param json body schemas.ModifyHole true "json"
// @Success 200 {object} Hole
// @Failure 404 {object} utils.MessageModel
func ModifyHole(c *fiber.Ctx) error {
	return nil
}

// DeleteHole
// @Summary Delete A Hole
// @Description Hide a hole, but visible to admins. This may affect many floors, DO NOT ABUSE!!!
// @Tags Hole
// @Produce application/json
// @Router /holes/{id} [delete]
// @Param id path int true "id"
// @Success 204
// @Failure 404 {object} utils.MessageModel
func DeleteHole(c *fiber.Ctx) error {
	return nil
}
