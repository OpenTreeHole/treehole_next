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

const MAX_GROUP_PER_USER = 10

type FavoriteGroups []FavoriteGroup

func (FavoriteGroup) TableName() string {
	return "favorite_groups"
}

func UserGetFavoriteGroups(tx *gorm.DB, userID int) (favoriteGroups FavoriteGroups, err error) {
	err = tx.Transaction(func(tx *gorm.DB) error {
		return tx.Where("user_id = ? and deleted = false", userID).Find(&favoriteGroups).Error
	})
	return
}

func DeleteUserFavoriteGroup(tx *gorm.DB, userID int, groupID int) (err error) {
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		return tx.Where("user_id = ? AND id = ?", userID, groupID).Updates(FavoriteGroup{Deleted: true}).Error
	})
}

func AddUserFavoriteGroup(tx *gorm.DB, userID int, name string) (err error) {
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		var groupID int
		err = tx.Raw("SELECT MAX(id) + 1 AS max_id FROM favorite_groups where user_id = ? and deleted = false", userID).Scan(&groupID).Error
		if err != nil {
			return err
		}
		if groupID >= MAX_GROUP_PER_USER {
			err = tx.Raw("SELECT id FROM favorite_groups where user_id = ? and deleted = true ORDER BY id LIMIT 1", userID).Take(&groupID).Error
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Forbidden("收藏夹数量已达上限")
		}
		if err != nil {
			return err
		}

		return tx.Create(&FavoriteGroup{
			UserID:   userID,
			Name:     name,
			ID:       groupID,
			CreateAt: time.Now(),
		}).Error
	})
}

func ModifyUserFavoriteGroup(tx *gorm.DB, userID int, groupID int, name string) (err error) {
	return tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		return tx.Where("user_id = ? AND id = ?", userID, groupID).Updates(FavoriteGroup{Name: name, UpdateAt: time.Now()}).Error
	})
}
