package favourite

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

// ListFavorites
//
//	@Summary	List User's Favorites
//	@Tags		Favorite
//	@Produce	application/json
//	@Router		/user/favorites [get]
//	@Param		object	query		ListModel	false	"query"
//	@Success	200		{object}	models.Map
//	@Success	200		{array}		models.Hole
func ListFavorites(c *fiber.Ctx) error {
	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	query, err := ValidateQuery[ListModel](c)
	if err != nil {
		return err
	}

	if query.Plain {
		// get favorite ids
		data, err := UserGetFavoriteData(DB, userID)
		if err != nil {
			return err
		}
		return c.JSON(Map{"data": data})
	} else {
		// get order
		var order string
		switch query.Order {
		case "id":
			order = "hole.id desc"
		case "time_created":
			order = "user_favorites.created_at desc, hole.id desc"
		case "hole_time_updated":
			order = "hole.updated_at desc"
		}

		// get favorites
		holes := make(Holes, 0)
		err = DB.
			Joins("JOIN user_favorites ON user_favorites.hole_id = hole.id AND user_favorites.user_id = ?", userID).
			Order(order).Find(&holes).Error
		if err != nil {
			return err
		}
		return Serialize(c, &holes)
	}
}

// AddFavorite
//
//	@Summary	Add A Favorite
//	@Tags		Favorite
//	@Accept		application/json
//	@Produce	application/json
//	@Router		/user/favorites [post]
//	@Param		json	body		AddModel	true	"json"
//	@Success	201		{object}	Response
//	@Success	200		{object}	Response
func AddFavorite(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[AddModel](c)
	if err != nil {
		return err
	}

	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	var data []int

	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		// add favorite
		err = AddUserFavourite(tx, userID, body.HoleID)
		if err != nil {
			return err
		}

		// create response
		data, err = UserGetFavoriteData(tx, userID)
		return err
	})
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&Response{
		Message: "收藏成功",
		Data:    data,
	})
}

// ModifyFavorite
//
//	@Summary	Modify User's Favorites
//	@Tags		Favorite
//	@Produce	application/json
//	@Router		/user/favorites [put]
//	@Param		json	body		ModifyModel	true	"json"
//	@Success	200		{object}	Response
//	@Failure	404		{object}	Response
func ModifyFavorite(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[ModifyModel](c)
	if err != nil {
		return err
	}

	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// modify favorite
	err = ModifyUserFavourite(DB, userID, body.HoleIDs)
	if err != nil {
		return err
	}

	// create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&Response{
		Message: "修改成功",
		Data:    data,
	})
}

// DeleteFavorite
//
//	@Summary	Delete A Favorite
//	@Tags		Favorite
//	@Produce	application/json
//	@Router		/user/favorites [delete]
//	@Param		json	body		DeleteModel	true	"json"
//	@Success	200		{object}	Response
//	@Failure	404		{object}	Response
func DeleteFavorite(c *fiber.Ctx) error {
	// validate body
	body, err := ValidateBody[DeleteModel](c)
	if err != nil {
		return err
	}

	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// delete favorite
	err = DB.Delete(UserFavorite{UserID: userID, HoleID: body.HoleID}).Error
	if err != nil {
		return err
	}

	// create response
	data, err := UserGetFavoriteData(DB, userID)
	if err != nil {
		return err
	}

	return c.JSON(&Response{
		Message: "删除成功",
		Data:    data,
	})
}
