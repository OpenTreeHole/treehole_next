package models

import (
	"errors"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"time"
)

// Punishment
// a record of user punishment
// when a record created, it can't be modified if other admins punish this user on the same floor
// whether a user is banned to post on one division based on the latest / max(id) record
// if admin want to modify punishment duration, manually modify the latest record of this user in database
// admin can be granted update privilege on SQL view of this table
type Punishment struct {
	ID int `json:"id" gorm:"primaryKey"`

	// time when this punishment creates
	CreatedAt time.Time `json:"created_at" gorm:"not null"`

	// start from end_time of previous punishment (punishment accumulation of different floors)
	// if no previous punishment or previous punishment end time less than time.Now() (synced), set start time time.Now()
	StartTime time.Time `json:"start_time" gorm:"not null"`

	// end_time of this punishment
	EndTime time.Time `json:"end_time" gorm:"not null"`

	Duration time.Duration `json:"duration" gorm:"not null"`

	// user punished
	UserID int `json:"user_id" gorm:"not null;index:idx_user_div,priority:1"`

	// admin user_id who made this punish
	MadeBy int `json:"made_by"`

	// punished because of this floor
	FloorID int `json:"floor_id" gorm:"not null;uniqueIndex:idx_user_floor,priority:2"`

	Floor *Floor `json:"floor"` // foreign key

	DivisionID int `json:"division_id" gorm:"not null"`

	Division *Division `json:"division"` // foreign key

	// reason
	Reason string `json:"reason" gorm:"size:128"`
}

type Punishments []*Punishment

func (punishment *Punishment) Create() (*User, error) {
	var user User

	err := DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, punishment.UserID).Error
		if err != nil {
			return err
		}

		var floorPunishment Punishment
		err = tx.Where("user_id = ? and floor_id = ?", user.ID, punishment.FloorID).Take(&floorPunishment).Error
		if err == nil {
			return common.Forbidden("该用户已被禁言")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var lastPunishment Punishment
		err = tx.Where("user_id = ?", user.ID).Last(&lastPunishment).Error
		if err == nil {
			if lastPunishment.EndTime.Before(time.Now()) {
				punishment.StartTime = time.Now()
			} else {
				punishment.StartTime = lastPunishment.EndTime
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			punishment.StartTime = time.Now()
		} else {
			return err
		}

		punishment.EndTime = punishment.StartTime.Add(punishment.Duration)
		user.BanDivision[punishment.DivisionID] = &punishment.EndTime
		user.OffenceCount += 1

		err = tx.Create(&punishment).Error
		if err != nil {
			return err
		}

		err = tx.Select("BanDivision", "OffenceCount").Save(&user).Error
		if err != nil {
			return err
		}

		return nil
	})
	return &user, err
}
