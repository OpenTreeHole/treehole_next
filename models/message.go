package models

// Should be same as message in notification project

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Messages []Message

type Message struct {
	ID          int         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time   `json:"time_created"`
	UpdatedAt   time.Time   `json:"time_updated"`
	Title       string      `json:"message" gorm:"size:32;not null"`
	Description string      `json:"description" gorm:"size:64;not null"`
	Data        any         `json:"data" gorm:"serializer:json" `
	Type        MessageType `json:"code" gorm:"size:16;not null"`
	URL         string      `json:"url" gorm:"size:64;default:'';not null"`
	Recipients  []int       `json:"-" gorm:"-:all" `
	MessageID   int         `json:"message_id" gorm:"-:all"`       // 兼容旧版 id
	HasRead     bool        `json:"has_read" gorm:"default:false"` // 兼容旧版, 永远为false，以MessageUser的HasRead为准
	Users       Users       `json:"-" gorm:"many2many:message_user;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type MessageUser struct {
	MessageID int  `json:"message_id" gorm:"primaryKey"`
	UserID    int  `json:"user_id" gorm:"primaryKey"`
	HasRead   bool `json:"has_read" gorm:"default:false"` // 兼容旧版
}

type MessageType string

const (
	MessageTypeFavorite    MessageType = "favorite"
	MessageTypeReply       MessageType = "reply"
	MessageTypeMention     MessageType = "mention"
	MessageTypeModify      MessageType = "modify" // including fold and delete
	MessageTypePermission  MessageType = "permission"
	MessageTypeReport      MessageType = "report"
	MessageTypeReportDealt MessageType = "report_dealt"
	MessageTypeMail        MessageType = "mail"
)

func (messages Messages) Preprocess(c *fiber.Ctx) error {
	for i := 0; i < len(messages); i++ {
		err := messages[i].Preprocess(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (message *Message) Preprocess(_ *fiber.Ctx) error {
	message.MessageID = message.ID
	return nil
}

func (message *Message) AfterCreate(tx *gorm.DB) (err error) {
	mapping := make([]MessageUser, len(message.Recipients))
	for i, userID := range message.Recipients {
		mapping[i] = MessageUser{
			MessageID: message.ID,
			UserID:    userID,
		}
	}
	return tx.Create(&mapping).Error
}
