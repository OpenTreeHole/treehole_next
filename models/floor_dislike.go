package models

import "gorm.io/gorm"

type FloorDislike struct {
	FloorID int `gorm:"primaryKey"`
	UserID  int `gorm:"primaryKey"`
}

func HasFloorDislike(tx *gorm.DB, floorID, userID int) (bool, error) {
	var floorDislike FloorDislike
	err := tx.Where("floor_id = ? and user_id = ?", floorID, userID).Take(&floorDislike).Error
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
