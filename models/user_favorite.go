package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserFavorite struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	HoleID    int       `json:"hole_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created"`
}

type UserFavorites []UserFavorite

func (UserFavorite) TableName() string {
	return "user_favorites"
}

//func makeUserFavorites(userID int, holeIDs []int) UserFavorites {
//	userFavorites := make(UserFavorites, 0, len(holeIDs))
//	for _, holeID := range holeIDs {
//		userFavorites = append(userFavorites, UserFavorite{
//			UserID: userID,
//			HoleID: holeID,
//		})
//	}
//	return userFavorites
//}

func ModifyUserFavourite(_ *gorm.DB, _ int, holeIDs []int) error {

	if len(holeIDs) == 0 {
		return nil
	}

	return nil
}

func AddUserFavourite(tx *gorm.DB, userID int, holeID int) error {
	return tx.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(Map{"created_at": time.Now()}),
	}).Create(&UserFavorite{
		UserID: userID,
		HoleID: holeID}).Error
}

func UserGetFavoriteData(tx *gorm.DB, userID int) ([]int, error) {
	data := make([]int, 0, 10)
	err := tx.Raw("SELECT hole_id FROM user_favorites WHERE user_id = ?", userID).Scan(&data).Error
	return data, err
}
