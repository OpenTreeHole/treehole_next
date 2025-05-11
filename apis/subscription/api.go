package subscription

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	. "treehole_next/models"
	. "treehole_next/utils"
)

// ListSubscriptions
//
// @Summary List User's Subscriptions
// @Tags Subscription
// @Produce application/json
// @Router /users/subscriptions [get]
// @Param object query ListModel false "query"
// @Success 200 {object} models.Map
// @Success 200 {array} models.Hole
func ListSubscriptions(c *fiber.Ctx) error {
	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	var query ListModel
	err = common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	if query.Plain {
		data, err := UserGetSubscriptionData(DB, userID)
		if err != nil {
			return err
		}
		return c.JSON(Map{"data": data})
	} else {
		holes := make(Holes, 0)
		err := DB.
			Joins("JOIN user_subscription ON user_subscription.hole_id = hole.id AND user_subscription.user_id = ?", userID).
			Order("user_subscription.created_at desc").Find(&holes).Error
		if err != nil {
			return err
		}
		return Serialize(c, &holes)
	}
}

// AddSubscription
//
// @Summary Add A Subscription
// @Tags Subscription
// @Accept application/json
// @Produce application/json
// @Router /users/subscriptions [post]
// @Param json body AddModel true "json"
// @Success 201 {object} Response
func AddSubscription(c *fiber.Ctx) error {
	// validate body
	var body AddModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	var data []int

	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		// add favorites
		err = AddUserSubscription(tx, userID, body.HoleID)
		if err != nil {
			return err
		}

		// create response
		data, err = UserGetSubscriptionData(tx, userID)
		return err
	})
	if err != nil {
		return err
	}

	return c.Status(201).JSON(&Response{
		Message: "关注成功",
		Data:    data,
	})
}

// DeleteSubscription
//
// @Summary Delete A Subscription
// @Tags Subscription
// @Produce application/json
// @Router /users/subscription [delete]
// @Param json body DeleteModel true "json"
// @Success 200 {object} Response
// @Failure 404 {object} Response
func DeleteSubscription(c *fiber.Ctx) error {
	// validate body
	var body DeleteModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get userID
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}
	var data []int

	err = DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		// 删除订阅并更新计数
		if err := RemoveUserSubscription(tx, userID, body.HoleID); err != nil {
			return err
		}

		var err error
		data, err = UserGetSubscriptionData(tx, userID)
		return err
	})

	return c.JSON(&Response{
		Message: "删除成功",
		Data:    data,
	})
}
