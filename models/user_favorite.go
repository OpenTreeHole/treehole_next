package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"time"
	"treehole_next/utils"
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

func ModifyUserFavourite(tx *gorm.DB, userID int, holeIDs []int) error {
	if len(holeIDs) == 0 {
		return nil
	}
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		var oldHoleIDs []int
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&UserFavorite{}).Select("hole_id").Scan(&oldHoleIDs).Error
		if err != nil {
			return err
		}

		// remove user_favorite that not in holeIDs
		var removingHoleIDMapping = make(map[int]bool)
		for _, holeID := range oldHoleIDs {
			removingHoleIDMapping[holeID] = true
		}
		for _, holeID := range holeIDs {
			if removingHoleIDMapping[holeID] {
				delete(removingHoleIDMapping, holeID)
			}
		}
		removingHoleIDs := utils.Keys(removingHoleIDMapping)
		if len(removingHoleIDs) > 0 {
			deleteUserFavorite := make(UserFavorites, 0)
			for _, holeID := range removingHoleIDs {
				deleteUserFavorite = append(deleteUserFavorite, UserFavorite{UserID: userID, HoleID: holeID})
			}
			err = tx.Delete(&deleteUserFavorite).Error
			if err != nil {
				return err
			}
		}

		// insert user_favorite that not in oldHoleIDs
		var newHoleIDMapping = make(map[int]bool)
		for _, holeID := range holeIDs {
			newHoleIDMapping[holeID] = true
		}
		for _, holeID := range oldHoleIDs {
			if newHoleIDMapping[holeID] {
				delete(newHoleIDMapping, holeID)
			}
		}
		newHoleIDs := utils.Keys(newHoleIDMapping)
		if len(newHoleIDs) > 0 {
			insertUserFavorite := make(UserFavorites, 0)
			for _, holeID := range newHoleIDs {
				insertUserFavorite = append(insertUserFavorite, UserFavorite{UserID: userID, HoleID: holeID})
			}
			err = tx.Create(&insertUserFavorite).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
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
