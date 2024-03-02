package models

import (
	"errors"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type FavoriteGroup struct {
	ID       int       `json:"id" gorm:"primaryKey"`
	UserID   int       `json:"user_id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"not null;size:64" default:"默认"`
	CreateAt time.Time `json:"time_created"`
	UpdateAt time.Time `json:"time_updated"`
	Deleted  bool      `json:"deleted" gorm:"default:false"`
	Number   int       `json:"number" gorm:"default:0"`
}

const MaxGroupPerUser = 10

type FavoriteGroups []FavoriteGroup

func (FavoriteGroup) TableName() string {
	return "favorite_groups"
}

func UserGetFavoriteGroups(tx *gorm.DB, userID int) (favoriteGroups FavoriteGroups, err error) {
	err = tx.Where("user_id = ? and deleted = false", userID).Find(&favoriteGroups).Error
	return
}

func DeleteUserFavoriteGroup(tx *gorm.DB, userID int, groupID int) (err error) {
	if groupID == 0 {
		return common.Forbidden("默认收藏夹不可删除")
	}
	err = tx.Clauses(dbresolver.Write).Where("user_id = ? AND id = ?", userID, groupID).Updates(FavoriteGroup{Deleted: true}).Error
	if err != nil {
		return err
	}
	err = tx.Model(&UserFavorite{}).Where("user_id = ? AND favorite_group_id = ?", userID, groupID).Delete(&UserFavorite{}).Error
	if err != nil {
		return err
	}
	return tx.Model(&User{}).Where("id = ?", userID).Update("favorite_group_count", gorm.Expr("favorite_group_count - 1")).Error
}

func AddUserFavoriteGroup(tx *gorm.DB, userID int, name string) (err error) {
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		var groupID int
		err = tx.Model(&FavoriteGroup{}).Select("MAX(id) + 1 AS max_id").Where("user_id = ? and deleted = false", userID).
			Take(&groupID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			groupID = 0
			err = nil
		}
		if err != nil {
			return err
		}
		if groupID >= MaxGroupPerUser {
			err = tx.Model(&FavoriteGroup{}).Where("user_id = ? and deleted = true", userID).Order("id").Limit(1).Take(&groupID).Error
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Forbidden("收藏夹数量已达上限")
		}
		if err != nil {
			return err
		}

		err = tx.Create(&FavoriteGroup{
			UserID:   userID,
			Name:     name,
			ID:       groupID,
			CreateAt: time.Now(),
		}).Error
		if err != nil {
			return err
		}
		return tx.Model(&User{}).Where("id = ?", userID).Update("favorite_group_count", gorm.Expr("favorite_group_count + 1")).Error
	})
}

func ModifyUserFavoriteGroup(tx *gorm.DB, userID int, groupID int, name string) (err error) {
	return tx.Clauses(dbresolver.Write).Where("user_id = ? AND id = ?", userID, groupID).
		Updates(FavoriteGroup{Name: name, UpdateAt: time.Now()}).Error
}
