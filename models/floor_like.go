package models

type FloorLike struct {
	FloorID  int  `json:"floor_id" gorm:"primaryKey"`
	UserID   int  `json:"user_id" gorm:"primaryKey"`
	LikeData int8 `json:"like_data"`
}
