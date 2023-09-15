package models

import (
	"errors"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"time"
)

type ReportPunishment struct {
	ID int `json:"id" gorm:"primaryKey"`

	// time when this report punishment creates
	CreatedAt time.Time `json:"created_at"`

	// time when this report punishment revoked
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	// start from end_time of previous punishment (punishment accumulation of different floors)
	// if no previous punishment or previous punishment end time less than time.Now() (synced), set start time time.Now()
	StartTime time.Time `json:"start_time" gorm:"not null"`

	// end_time of this report punishment
	EndTime time.Time `json:"end_time" gorm:"not null"`

	Duration *time.Duration `json:"duration"`

	// user punished
	UserID int `json:"user_id" gorm:"not null;index"`

	// admin user_id who made this punish
	MadeBy int `json:"made_by,omitempty"`

	// punished because of this report
	ReportId int `json:"report_id" gorm:"uniqueIndex"`

	Report *Report `json:"report,omitempty"` // foreign key

	// reason
	Reason string `json:"reason" gorm:"size:128"`
}

type ReportPunishments []*ReportPunishment

func (reportPunishment *ReportPunishment) Create() (*User, error) {
	var user User

	err := DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, reportPunishment.UserID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// if the user has been banned from this report
		var punishmentRecord ReportPunishment
		err = tx.Where("user_id = ? and report_id = ?", user.ID, reportPunishment.ReportId).Take(&punishmentRecord).Error
		if err == nil {
			return common.Forbidden("该用户已被限制使用举报功能")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var lastPunishment ReportPunishment
		err = tx.Where("user_id = ?", user.ID).Last(&lastPunishment).Error
		if err == nil {
			if lastPunishment.EndTime.Before(time.Now()) {
				reportPunishment.StartTime = time.Now()
			} else {
				reportPunishment.StartTime = lastPunishment.EndTime
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			reportPunishment.StartTime = time.Now()
		} else {
			return err
		}

		reportPunishment.EndTime = reportPunishment.StartTime.Add(*reportPunishment.Duration)

		user.BanReport = &reportPunishment.EndTime
		user.BanReportCount += 1

		err = tx.Create(&reportPunishment).Error
		if err != nil {
			return err
		}

		err = tx.Select("BanReport", "BanReportCount").Save(&user).Error
		if err != nil {
			return err
		}

		return nil
	})
	return &user, err
}
