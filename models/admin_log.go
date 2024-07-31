package models

import (
	"github.com/rs/zerolog/log"
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
	AdminLogTypeHideHole        AdminLogType = "hide_hole"
	AdminLogTypeTag             AdminLogType = "edit_tag"
	AdminLogTypeDivision        AdminLogType = "edit_division"
	AdminLogTypeMessage         AdminLogType = "send_message"
	AdminLogTypeDeleteReport    AdminLogType = "delete_report"
	AdminLogTypeChangeSensitive AdminLogType = "change_sensitive"
)

// CreateAdminLog
// save admin edit log for audit purpose
func CreateAdminLog(tx *gorm.DB, logType AdminLogType, userID int, data any) {
	adminLog := AdminLog{
		Type:   logType,
		UserID: userID,
		Data:   data,
	}
	err := tx.Create(&adminLog).Error // omit error
	if err != nil {
		log.Error().Err(err).Msg("failed to create admin log")
	}
}
