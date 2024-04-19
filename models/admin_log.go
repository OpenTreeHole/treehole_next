package models

import (
	"gorm.io/gorm"
	"time"
)

type AdminLog struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	Type      AdminLogType `gorm:"size:16;not null"`
	UserID    int          `gorm:"not null"`
	Data      any          `gorm:"serializer:json"`
}

type AdminLogType string

const (
	AdminLogTypeHole            AdminLogType = "edit_hole"
	AdminLogTypeDivision        AdminLogType = "edit_division"
	AdminLogTypeMessage         AdminLogType = "send_message"
	AdminLogTypeChangeSensitive AdminLogType = "change_sensitive"
)

// CreateAdminLog
// save admin edit log for audit purpose
func CreateAdminLog(tx *gorm.DB, logType AdminLogType, userID int, data any) {
	log := AdminLog{
		Type:   logType,
		UserID: userID,
		Data:   data,
	}
	tx.Create(&log) // omit error
}
