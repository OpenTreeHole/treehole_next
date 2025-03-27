package message

import (
	"github.com/opentreehole/go-common"

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
	var query ListModel
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	userID, err := common.GetUserID(c)
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
			userID,
		).Scan(&messages)
	} else {
		DB.Raw(`
			SELECT message.*,message_user.has_read FROM message
			INNER JOIN message_user
			WHERE message.id = message_user.message_id and message_user.user_id = ?
			ORDER BY updated_at DESC`,
			userID,
		).Scan(&messages)
	}

	return Serialize(c, &messages)
}

// SendMail
// @Summary Send a Mail
// @Description Send to multiple recipients and save to db, admin only.
// @Tags Message
// @Produce application/json
// @Param json body CreateModel true "json"
// @Router /messages [post]
// @Success 201 {object} Message
func SendMail(c *fiber.Ctx) error {
	var body CreateModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get user
	user, err := GetCurrLoginUser(c)
	if err != nil {
		return err
	}

	// permission
	if !user.IsAdmin {
		return common.Forbidden()
	}

	// construct mail
	mail := Notification{
		Description: body.Description,
		Recipients:  body.Recipients,
		Data:        Map{},
		Title:       "您有一封站内信",
		Type:        MessageTypeMail,
		URL:         "/api/messages",
	}

	// send
	message, err := mail.Send()
	if err != nil {
		return err
	}

	CreateAdminLog(DB, AdminLogTypeMessage, user.ID, body)

	return Serialize(c.Status(201), &message)
}

// ClearMessages
// @Summary Clear Messages of a User
// @Tags Message
// @Produce application/json
// @Router /messages/clear [post]
// @Success 204
func ClearMessages(c *fiber.Ctx) error {
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	result := DB.Exec(
		"UPDATE message_user SET has_read = true WHERE user_id = ?",
		userID,
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
// @Router /messages/_webvpn [patch]
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
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	id, _ := c.ParamsInt("id")
	result := DB.Exec(
		"UPDATE message_user SET has_read = true WHERE user_id = ?  AND message_id = ?",
		userID, id,
	)
	if result.Error != nil {
		return result.Error
	}
	return c.Status(204).JSON(nil)
}
