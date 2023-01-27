package models

type AnonynameMapping struct {
	HoleID    int    `json:"hole_id" gorm:"primaryKey"`
	UserID    int    `json:"user_id" gorm:"primaryKey"`
	Anonyname string `json:"anonyname" gorm:"size:32"`
}
