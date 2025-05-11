package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
)

type UserSubscription struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	HoleID    int       `json:"hole_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created"`
}

type UserSubscriptions []UserSubscription

func (UserSubscription) TableName() string {
	return "user_subscription"
}

func UserGetSubscriptionData(tx *gorm.DB, userID int) ([]int, error) {
	data := make([]int, 0, 10)
	err := tx.Clauses(dbresolver.Write).Raw("SELECT hole_id FROM user_subscription WHERE user_id = ? ORDER BY created_at", userID).Scan(&data).Error
	return data, err
}

func AddUserSubscription(tx *gorm.DB, userID int, holeID int) error {

	// 检查是否已存在订阅关系
	var exists int64
	if err := tx.Model(&UserSubscription{}).Where("user_id = ? AND hole_id = ?", userID, holeID).Count(&exists).Error; err != nil {
		return err
	}

	err := tx.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"created_at": time.Now()}),
	}).Create(&UserSubscription{
		UserID: userID,
		HoleID: holeID,
	}).Error

	if err != nil {
		return err
	}

	if exists == 0 {
		return tx.Model(&Hole{}).Where("id = ?", holeID).
			UpdateColumn("subscription_count", gorm.Expr("subscription_count + ?", 1)).Error
	}

	return nil
}

func RemoveUserSubscription(tx *gorm.DB, userID int, holeID int) error {
	// 检查记录是否存在
	var exists int64
	if err := tx.Model(&UserSubscription{}).Where("user_id = ? AND hole_id = ?", userID, holeID).Count(&exists).Error; err != nil {
		return err
	}

	// 只有当记录存在时才执行删除和计数更新
	if exists > 0 {
		// 删除订阅
		if err := tx.Where("user_id = ? AND hole_id = ?", userID, holeID).Delete(&UserSubscription{}).Error; err != nil {
			return err
		}

		// 更新订阅计数
		return tx.Model(&Hole{}).Where("id = ?", holeID).
			UpdateColumn("subscription_count", gorm.Expr("subscription_count - ?", 1)).Error
	}

	return nil
}
