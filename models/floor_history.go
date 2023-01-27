package models

import "time"

type FloorHistory struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created"`
	UpdatedAt time.Time `json:"time_updated"`
	Content   string    `json:"content"`
	Reason    string    `json:"reason"`
	FloorID   int       `json:"floor_id"`
	UserID    int       `json:"user_id"` // The one who modified the floor
}
