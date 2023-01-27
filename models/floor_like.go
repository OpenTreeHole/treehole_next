package models

import "gorm.io/gorm"

type FloorLike struct {
	FloorID  int  `json:"floor_id" gorm:"primaryKey"`
	UserID   int  `json:"user_id" gorm:"primaryKey"`
	LikeData int8 `json:"like_data"`
}

func HasFloorLike(tx *gorm.DB, floorID, userID int) (bool, error) {
	var floorLike FloorLike
	err := tx.Where("floor_id = ? and user_id = ?", floorID, userID).Take(&floorLike).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
