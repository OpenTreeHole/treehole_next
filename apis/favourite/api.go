package favourite

import "github.com/gofiber/fiber/v2"

// ListFavorites
// @Summary List User's Favorites
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [get]
// @Success 200 {array} models.Hole
func ListFavorites(c *fiber.Ctx) error {
	return nil
}

// AddFavorite
// @Summary Add A Favorite
// @Tags Favorite
// @Accept application/json
// @Produce application/json
// @Router /user/favorites [post]
// @Param json body AddModel true "json"
// @Success 201 {object} models.MessageModel
// @Success 200 {object} models.MessageModel
func AddFavorite(c *fiber.Ctx) error {
	return nil
}

// ModifyFavorite
// @Summary Modify User's Favorites
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [put]
// @Param json body ModifyModel true "json"
// @Success 200 {object} models.MessageModel
// @Failure 404 {object} models.MessageModel
func ModifyFavorite(c *fiber.Ctx) error {
	return nil
}

// DeleteFavorite
// @Summary Delete A Favorite
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [delete]
// @Param json body DeleteModel true "json"
// @Success 204
// @Failure 404 {object} models.MessageModel
func DeleteFavorite(c *fiber.Ctx) error {
	return nil
}
