package message

import (
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"

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

	q := DB.Model(&Message{}).Select("message.*,message_user.has_read").Joins("JOIN message_user ON message.id=message_user.message_id").Where("message_user.user_id = ?", userID)
	if query.NotRead {
		q = q.Where("message_user.has_read = ?", false)
	}
	uid, err := FilterUserID(c)
	if err != nil {
		return err
	}
	if uid != 0 {
		// Older notifications stored their canonical source only in the URL.
		// Resolve those links in bulk, then apply the same policy as newer messages.
		q = q.Joins(`LEFT JOIN report legacy_report ON message.related_floor_id IS NULL AND
 message.url LIKE '/api/reports/%' AND legacy_report.id=CAST(SUBSTR(message.url,14) AS SIGNED)`)
		floorSource := `COALESCE(message.related_floor_id,legacy_report.floor_id,
 CASE WHEN message.url LIKE '/api/floors/%' THEN CAST(SUBSTR(message.url,13) AS SIGNED) END)`
		holeSource := `COALESCE(message.related_hole_id,
 CASE WHEN message.url LIKE '/api/holes/%' THEN CAST(SUBSTR(message.url,12) AS SIGNED) END)`
		q = q.Where("(message.type = ? OR "+floorSource+" IS NOT NULL OR "+holeSource+" IS NOT NULL)", MessageTypeMail).
			Where("("+floorSource+" IS NULL OR "+floorSource+" IN (?))", VisibleFloorIDs(DB, uid)).
			Where("("+holeSource+" IS NULL OR "+holeSource+" IN (?))", VisiblePublicHoleIDs(DB, uid))
	}
	if err = q.Order("message.updated_at DESC").Scan(&messages).Error; err != nil {
		return err
	}
	if uid != 0 {
		for i := range messages {
			messages[i].Data = nil
		}
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

func deleteMessagesByCondition(db *gorm.DB, condition string, value int) error {
	var messageIDs []int
	if err := db.Model(&Message{}).Where(condition, value).Pluck("id", &messageIDs).Error; err != nil {
		return err
	}
	if len(messageIDs) == 0 {
		return nil
	}
	if err := db.Where("message_id IN ?", messageIDs).Delete(&MessageUser{}).Error; err != nil {
		return err
	}
	return db.Where("id IN ?", messageIDs).Delete(&Message{}).Error
}

func DeleteMessageByRelatedFloorID(db *gorm.DB, floorID int) error {
	return deleteMessagesByCondition(db, "related_floor_id = ?", floorID)
}

func DeleteMessageByRelatedHoleID(db *gorm.DB, holeID int) error {
	return deleteMessagesByCondition(db, "related_hole_id = ?", holeID)
}
