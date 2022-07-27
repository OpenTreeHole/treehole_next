package favourite

import (
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

// ListFavorites
// @Summary List User's Favorites
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [get]
// @Success 200 {array} models.Hole
func ListFavorites(c *fiber.Ctx) error {
	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// get favorites
	var holes Holes
	sql := `SELECT * FROM hole 
	JOIN user_favorites 
	ON user_favorites.hole_id = hole.id 
	AND user_favorites.user_id = ?`

	DB.Raw(sql, userID).Scan(&holes)
	return Serialize(c, &holes)
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
	// validate body
	var body AddModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// add favorite
	err = UserCreateFavourite(c, false, userID, []int{body.HoleID})
	if err != nil {
		return err
	}
	return c.JSON(MessageModel{Message: "收藏成功"})
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
	// validate body
	var body ModifyModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// modify favorite
	err = UserCreateFavourite(c, true, userID, body.HoleIDs)
	if err != nil {
		return err
	}
	return c.JSON(MessageModel{Message: "修改成功"})
}

// DeleteFavorite
// @Summary Delete A Favorite
// @Tags Favorite
// @Produce application/json
// @Router /user/favorites [delete]
// @Param json body DeleteModel true "json"
// @Success 200
// @Failure 404 {object} models.MessageModel
func DeleteFavorite(c *fiber.Ctx) error {
	// validate body
	var body DeleteModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// modify favorite
	err = UserDeleteFavorite(userID, []int{body.HoleID})
	if err != nil {
		return err
	}
	return c.JSON(MessageModel{Message: "删除成功"})
}
