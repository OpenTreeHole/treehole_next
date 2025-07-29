package models

import (
	"github.com/opentreehole/go-common"
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

func IsFavoriteGroupExist(tx *gorm.DB, userID int, favoriteGroupID int) bool {
	var num int64
	tx.Model(&FavoriteGroup{}).Where("user_id = ? AND favorite_group_id = ? AND deleted = false", userID, favoriteGroupID).Count(&num)
	return num > 0
}

// ModifyUserFavorite only take effect in the same favorite_group
func ModifyUserFavorite(tx *gorm.DB, userID int, holeIDs []int, favoriteGroupID int) error {
	if len(holeIDs) == 0 {
		return nil
	}
	if !IsFavoriteGroupExist(tx, userID, favoriteGroupID) {
		return common.NotFound("收藏夹不存在")
	}
	if !IsHolesExist(tx, holeIDs) {
		return common.Forbidden("帖子不存在")
	}
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		var oldHoleIDs []int
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&UserFavorite{}).Where("user_id = ? AND favorite_group_id = ?", userID, favoriteGroupID).
			Pluck("hole_id", &oldHoleIDs).Error
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
				deleteUserFavorite = append(deleteUserFavorite, UserFavorite{UserID: userID, HoleID: holeID, FavoriteGroupID: favoriteGroupID})
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
				insertUserFavorite = append(insertUserFavorite, UserFavorite{UserID: userID, HoleID: holeID, FavoriteGroupID: favoriteGroupID})
			}
			err = tx.Create(&insertUserFavorite).Error
			if err != nil {
				return err
			}
		}
		return tx.Model(&FavoriteGroup{}).Where("user_id = ? AND favorite_group_id = ?", userID, favoriteGroupID).Update("count", len(holeIDs)).Error
	})
}

func AddUserFavorite(tx *gorm.DB, userID int, holeID int, favoriteGroupID int) error {
	if !IsFavoriteGroupExist(tx, userID, favoriteGroupID) {
		return common.NotFound("收藏夹不存在")
	}
	if !IsHolesExist(tx, []int{holeID}) {
		return common.NotFound("帖子不存在")
	}
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
	err = tx.Clauses(dbresolver.Write).Model(&FavoriteGroup{}).
		Where("user_id = ? AND favorite_group_id = ?", userID, favoriteGroupID).Update("count", gorm.Expr("count + 1")).Error
	if err != nil {
		return err
	}
	return tx.Model(&Hole{}).Where("id = ?", holeID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error

}

// UserGetFavoriteData get all favorite data of a user
func UserGetFavoriteData(tx *gorm.DB, userID int) ([]int, error) {
	data := make([]int, 0, 10)
	err := tx.Clauses(dbresolver.Write).Model(&UserFavorite{}).Where("user_id = ?", userID).Distinct().
		Pluck("hole_id", &data).Error
	return data, err
}

// UserGetFavoriteDataByFavoriteGroup get favorite data in specific favorite group
func UserGetFavoriteDataByFavoriteGroup(tx *gorm.DB, userID int, favoriteGroupID int) ([]int, error) {
	if !IsFavoriteGroupExist(tx, userID, favoriteGroupID) {
		return nil, common.NotFound("收藏夹不存在")
	}
	data := make([]int, 0, 10)
	err := tx.Clauses(dbresolver.Write).Model(&UserFavorite{}).
		Where("user_id = ? AND favorite_group_id = ?", userID, favoriteGroupID).Pluck("hole_id", &data).Error
	return data, err
}

// DeleteUserFavorite delete user favorite
// if user favorite hole only once, delete the hole
// otherwise, delete the favorite in the specific favorite group
func DeleteUserFavorite(tx *gorm.DB, userID int, holeID int, favoriteGroupID int) error {
	if !IsFavoriteGroupExist(tx, userID, favoriteGroupID) {
		return common.NotFound("收藏夹不存在")
	}
	if !IsHolesExist(tx, []int{holeID}) {
		return common.NotFound("帖子不存在")
	}

	// 检查记录是否存在
	var count int64
	if err := tx.Model(&UserFavorite{}).
		Where("user_id = ? AND hole_id = ? AND favorite_group_id = ?",
			userID, holeID, favoriteGroupID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		// 删除收藏记录
		if err := tx.Where("user_id = ? AND hole_id = ? AND favorite_group_id = ?",
			userID, holeID, favoriteGroupID).
			Delete(&UserFavorite{}).Error; err != nil {
			return err
		}

		// 更新收藏夹计数
		if err := tx.Model(&FavoriteGroup{}).
			Where("favorite_group_id = ? AND user_id = ?", favoriteGroupID, userID).
			UpdateColumn("count", gorm.Expr("count - ?", 1)).Error; err != nil {
			return err
		}

		// 更新帖子收藏计数
		if err := tx.Model(&Hole{}).Where("id = ?", holeID).
			UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}
	}

	return nil
}

// MoveUserFavorite move holes that are really in the fromFavoriteGroup
func MoveUserFavorite(tx *gorm.DB, userID int, holeIDs []int, fromFavoriteGroupID int, toFavoriteGroupID int) error {
	if fromFavoriteGroupID == toFavoriteGroupID {
		return nil
	}
	if len(holeIDs) == 0 {
		return nil
	}
	if !IsFavoriteGroupExist(tx, userID, fromFavoriteGroupID) || !IsFavoriteGroupExist(tx, userID, toFavoriteGroupID) {
		return common.NotFound("收藏夹不存在")
	}
	if !IsHolesExist(tx, holeIDs) {
		return common.Forbidden("帖子不存在")
	}
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		var oldHoleIDs []int
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&UserFavorite{}).Where("user_id = ? AND favorite_group_id = ?", userID, fromFavoriteGroupID).
			Pluck("hole_id", &oldHoleIDs).Error
		if err != nil {
			return err
		}

		// move user_favorite that in holeIDs
		var removingHoleIDMapping = make(map[int]bool)
		var removingHoleIDs []int
		for _, holeID := range oldHoleIDs {
			removingHoleIDMapping[holeID] = true
		}
		for _, holeID := range holeIDs {
			if removingHoleIDMapping[holeID] {
				removingHoleIDs = append(removingHoleIDs, holeID)
			}
		}
		if len(removingHoleIDs) > 0 {
			err = tx.Table("user_favorites").
				Where("user_id = ? AND favorite_group_id = ? AND hole_id IN ?", userID, fromFavoriteGroupID, removingHoleIDs).
				Updates(map[string]interface{}{"favorite_group_id": toFavoriteGroupID}).Error
			if err != nil {
				return err
			}
		}
		err = tx.Model(&FavoriteGroup{}).Where("user_id = ? AND favorite_group_id = ?", userID, fromFavoriteGroupID).Update("count", gorm.Expr("count - ?", len(removingHoleIDs))).Error
		if err != nil {
			return err
		}
		return tx.Model(&FavoriteGroup{}).Where("user_id = ? AND favorite_group_id = ?", userID, toFavoriteGroupID).Update("count", gorm.Expr("count + ?", len(removingHoleIDs))).Error
	})
}
