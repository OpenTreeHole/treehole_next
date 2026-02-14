package models

import (
	"errors"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type FavoriteGroup struct {
	FavoriteGroupID int       `json:"favorite_group_id" gorm:"primaryKey"`
	UserID          int       `json:"user_id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null;size:64" default:"默认"`
	CreatedAt       time.Time `json:"time_created"`
	UpdatedAt       time.Time `json:"time_updated"`
	Deleted         bool      `json:"deleted" gorm:"default:false"`
	Count           int       `json:"count" gorm:"default:0"`
}

const MaxGroupPerUser = 10

type FavoriteGroups []FavoriteGroup

func (FavoriteGroup) TableName() string {
	return "favorite_groups"
}

// make sure use this function in a transaction
func UserGetFavoriteGroups(tx *gorm.DB, userID int, order *string) (favoriteGroups FavoriteGroups, err error) {
	err = CheckDefaultFavoriteGroup(tx, userID)
	if err != nil {
		return
	}

	if order == nil {
		err = tx.Where("user_id = ? and deleted = false", userID).Find(&favoriteGroups).Error
	} else {
		err = tx.Where("user_id = ? and deleted = false", userID).Order(*order).Find(&favoriteGroups).Error
	}
	return
}

func DeleteUserFavoriteGroup(tx *gorm.DB, userID int, groupID int) (err error) {
	if groupID == 0 {
		return common.Forbidden("默认收藏夹不可删除")
	}
	err = tx.Model(&UserFavorite{}).Where("user_id = ? AND favorite_group_id = ?", userID, groupID).Take(&UserFavorite{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		return common.Forbidden("收藏夹中存在收藏内容，请先移除")
	}

	result := tx.Clauses(dbresolver.Write).Where("user_id = ? AND favorite_group_id = ?", userID, groupID).Updates(FavoriteGroup{Deleted: true})
	if result.Error != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return common.NotFound("收藏夹不存在")
	}
	err = tx.Model(&UserFavorite{}).Where("user_id = ? AND favorite_group_id = ?", userID, groupID).Delete(&UserFavorite{}).Error
	if err != nil {
		return err
	}
	return tx.Model(&User{}).Where("id = ?", userID).Update("favorite_group_count", gorm.Expr("favorite_group_count - 1")).Error
}

func CheckDefaultFavoriteGroup(tx *gorm.DB, userID int) (err error) {
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&FavoriteGroup{}).Where("user_id = ? AND favorite_group_id = 0", userID).Take(&FavoriteGroup{}).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			// insert default favorite group if not exists
			err = tx.Create(&FavoriteGroup{
				UserID:          userID,
				Name:            "默认收藏夹",
				FavoriteGroupID: 0,
				CreatedAt:       time.Now(),
			}).Error
			if err != nil {
				return err
			}
			return tx.Model(&User{}).Where("id = ?", userID).Update("favorite_group_count", gorm.Expr("favorite_group_count + 1")).Error
		}

		// default favorite group exists
		return nil
	})

}

func AddUserFavoriteGroup(tx *gorm.DB, userID int, name string) (err error) {
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		var groupID int
		err = tx.Model(&FavoriteGroup{}).Select("IFNULL(MAX(favorite_group_id), 0) AS max_id").Where("user_id = ? and deleted = false", userID).
			Take(&groupID).Error
		groupID++
		if err != nil {
			return err
		}
		if groupID >= MaxGroupPerUser {
			err = tx.Model(&FavoriteGroup{}).Where("user_id = ? and deleted = true", userID).Order("favorite_group_id").Limit(1).Take(&groupID).Error
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Forbidden("收藏夹数量已达上限")
		}
		if err != nil {
			return err
		}

		err = tx.Create(&FavoriteGroup{
			UserID:          userID,
			Name:            name,
			FavoriteGroupID: groupID,
			CreatedAt:       time.Now(),
		}).Error
		if err != nil {
			return err
		}
		return tx.Model(&User{}).Where("id = ?", userID).Update("favorite_group_count", gorm.Expr("favorite_group_count + 1")).Error
	})
}

func ModifyUserFavoriteGroup(tx *gorm.DB, userID int, groupID int, name string) (err error) {
	return tx.Clauses(dbresolver.Write).Where("user_id = ? AND favorite_group_id = ?", userID, groupID).
		Updates(FavoriteGroup{Name: name, UpdatedAt: time.Now()}).Error
}
