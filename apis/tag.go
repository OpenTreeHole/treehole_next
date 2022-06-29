package apis

import (
	"fmt"
	. "treehole_next/models"

	"github.com/gofiber/fiber/v2"
)

// ListTags
// @Summary List All Tags
// @Tags Tag
// @Produce application/json
// @Router /tags [get]
// @Success 200 {array} Tag
func ListTags(c *fiber.Ctx) error {
	fmt.Println(Tag{})
	return nil
}

// GetTag
// @Summary Get A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Tag
// @Failure 404 {object} schemas.MessageModel
func GetTag(c *fiber.Ctx) error {
	return nil
}

// CreateTag
// @Summary Create A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags [post]
// @Param json body schemas.CreateTag true "json"
// @Success 200 {object} Tag
// @Success 201 {object} Tag
func CreateTag(c *fiber.Ctx) error {
	return nil
}

// ModifyTag
// @Summary Modify A Tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [put]
// @Param id path int true "id"
// @Param json body schemas.ModifyTag true "json"
// @Success 200 {object} Tag
// @Failure 404 {object} schemas.MessageModel
func ModifyTag(c *fiber.Ctx) error {
	return nil
}

// DeleteTag
// @Summary Delete A Tag
// @Description Delete a tag and link all of its holes to another given tag
// @Tags Tag
// @Produce application/json
// @Router /tags/{id} [delete]
// @Param id path int true "id"
// @Param json body schemas.DeleteTag true "json"
// @Success 200 {object} Tag
// @Failure 404 {object} schemas.MessageModel
func DeleteTag(c *fiber.Ctx) error {
	return nil
}
