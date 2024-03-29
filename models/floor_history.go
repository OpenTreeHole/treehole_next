package models

import "time"

type FloorHistory struct {
	/// base info
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created"`
	UpdatedAt time.Time `json:"time_updated"`
	Content   string    `json:"content" gorm:"size:10000"`
	Reason    string    `json:"reason"`
	FloorID   int       `json:"floor_id"`
	// auto sensitive check
	IsSensitive bool `json:"is_sensitive"`

	// manual sensitive check
	IsActualSensitive *bool `json:"is_actual_sensitive"`
	// The one who modified the floor
	UserID int `json:"user_id"`
}

type FloorHistorySlice []*FloorHistory
