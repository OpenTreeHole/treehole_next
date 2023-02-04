// Should be same as message in notification project
package models

import (
	"database/sql/driver"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSON map[string]any

func (t JSON) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *JSON) Scan(input any) error {
	return json.Unmarshal(input.([]byte), t)
}

// GormDataType gorm common data type
func (JSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
//
//goland:noinspection GoUnusedParameter
func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

type Messages []Message

type Message struct {
	ID          int         `gorm:"primarykey" json:"id"`
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
	Users       Users       `json:"user" gorm:"many2many:message_user;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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

func (message *Message) Preprocess(c *fiber.Ctx) error {
	message.MessageID = message.ID
	return nil
}

func (m *Message) AfterCreate(tx *gorm.DB) (err error) {
	mapping := make([]MessageUser, len(m.Recipients))
	for i, userID := range m.Recipients {
		mapping[i] = MessageUser{
			MessageID: m.ID,
			UserID:    userID,
		}
	}
	return tx.Create(&mapping).Error
}
