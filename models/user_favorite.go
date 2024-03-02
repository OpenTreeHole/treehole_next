package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"treehole_next/utils"
)

type UserFavorite struct {
	UserID          int       `json:"user_id" gorm:"primaryKey"`
	FavoriteGroupID int       `json:"favorite_group_id" gorm:"primaryKey"`
	HoleID          int       `json:"hole_id" gorm:"primaryKey"`
	CreatedAt       time.Time `json:"time_created"`
}

type UserFavorites []UserFavorite

func (UserFavorite) TableName() string {
	return "user_favorites"
}

// ModifyUserFavorite only take effect in the same favorite_group
// todo
func ModifyUserFavorite(tx *gorm.DB, userID int, holeIDs []int, favoriteGroupID int) error {
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

func AddUserFavorite(tx *gorm.DB, userID int, holeID int, favoriteGroupID int) error {
	var err = tx.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(Map{"created_at": time.Now()}),
	}).Create(&UserFavorite{
		UserID:          userID,
		HoleID:          holeID,
		FavoriteGroupID: favoriteGroupID,
	}).Error
	if err != nil {
		return err
	}
	return tx.Clauses(dbresolver.Write).Model(&FavoriteGroup{}).
		Where("user_id = ? AND id = ?", userID, favoriteGroupID).Update("number", gorm.Expr("number + 1")).Error
}

// UserGetFavoriteData get all favorite data of a user
func UserGetFavoriteData(tx *gorm.DB, userID int) ([]int, error) {
	data := make([]int, 0, 10)
	err := tx.Clauses(dbresolver.Write).Raw("SELECT DISTINCT hole_id FROM user_favorites WHERE user_id = ?", userID).Scan(&data).Error
	return data, err
}

// UserGetFavoriteDataByFavoriteGroup get favorite data in specific favorite group
func UserGetFavoriteDataByFavoriteGroup(tx *gorm.DB, userID int, favoriteGroupID int) ([]int, error) {
	data := make([]int, 0, 10)
	err := tx.Clauses(dbresolver.Write).Raw("SELECT hole_id FROM user_favorites WHERE user_id = ? and favorite_group_id = ?", userID, favoriteGroupID).Scan(&data).Error
	return data, err
}

func DeleteUserFavorite(tx *gorm.DB, userID int, holeID int, favoriteGroupID int) error {
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		var num int64
		err := tx.Model(&UserFavorite{}).Where("user_id = ? AND hole_id = ?", userID, holeID).Count(&num).Error
		if err != nil {
			return err
		}
		if num == 1 {
			err = tx.Delete(&UserFavorite{UserID: userID, HoleID: holeID}).Error
			if err != nil {
				return err
			}
			return tx.Clauses(dbresolver.Write).Model(&FavoriteGroup{}).Where("user_id = ? AND id = ?", userID, 0).Update("number", gorm.Expr("number - 1")).Error
		}
		err = tx.Delete(&UserFavorite{UserID: userID, HoleID: holeID, FavoriteGroupID: favoriteGroupID}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(dbresolver.Write).Model(&FavoriteGroup{}).Where("user_id = ? AND id = ?", userID, favoriteGroupID).Update("number", gorm.Expr("number - 1")).Error
	})
}
