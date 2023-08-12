package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"time"
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
	return tx.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(Map{"created_at": time.Now()}),
	}).Create(&UserSubscription{
		UserID: userID,
		HoleID: holeID}).Error
}
