package message

import (
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
)

// ListMessages
// @Summary List Messages of a User
// @Tags Message
// @Produce application/json
// @Router /messages [get]
// @Success 200 {array} Message
// @Param object query ListModel false "query"
func ListMessages(c *fiber.Ctx) error {
	// swagger里面的query是错误的，应该使用not_read而不是notRead
	query, err := ValidateQuery[ListModel](c)
	if err != nil {
		return err
	}

	messages := Messages{}

	if query.NotRead {
		DB.Raw(`
			SELECT message.*,message_user.has_read FROM message
			INNER JOIN message_user 
			WHERE message.id = message_user.message_id and message_user.user_id = ? and message_user.has_read = false
			ORDER BY updated_at DESC`,
			c.Locals("userID").(int),
		).Scan(&messages)
	} else {
		DB.Raw(`
			SELECT message.*,message_user.has_read FROM message
			INNER JOIN message_user
			WHERE message.id = message_user.message_id and message_user.user_id = ?
			ORDER BY updated_at DESC`,
			c.Locals("userID").(int),
		).Scan(&messages)
	}

	return Serialize(c, &messages)
}

// ClearMessages
// @Summary Clear Messages of a User
// @Tags Message
// @Produce application/json
// @Router /messages/clear [post]
// @Success 204
func ClearMessages(c *fiber.Ctx) error {
	result := DB.Exec(
		"UPDATE message_user SET has_read = true WHERE user_id = ?",
		c.Locals("userID").(int),
	)
	if result.Error != nil {
		return result.Error
	}
	return c.Status(204).JSON(nil)
}

// ClearMessagesDeprecated
// @Summary Clear Messages Deprecated
// @Tags Message
// @Produce application/json
// @Router /messages [put]
// @Success 204
func ClearMessagesDeprecated(c *fiber.Ctx) error {
	return ClearMessages(c)
}

// DeleteMessage
// @Summary Delete a message of a user
// @Tags Message
// @Produce application/json
// @Router /messages/{id} [delete]
// @Param id path int true "message id"
// @Success 204
func DeleteMessage(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	result := DB.Exec(
		"UPDATE message_user SET has_read = true WHERE user_id = ?  AND message_id = ?",
		c.Locals("userID").(int), id,
	)
	if result.Error != nil {
		return result.Error
	}
	return c.Status(204).JSON(nil)
}
