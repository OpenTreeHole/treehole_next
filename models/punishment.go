package models

import (
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
	CreateAt time.Time `json:"create_at" gorm:"not null"`

	// start from end_time of previous punishment (punishment accumulation of different floors)
	// if no previous punishment or previous punishment end time less than time.Now() (synced), set start time time.Now()
	StartTime time.Time `json:"start_time" gorm:"not null"`

	// end_time of this punishment
	EndTime time.Time `json:"end_time" gorm:"not null"`

	// user punished
	UserID int `json:"user_id" gorm:"index:idx_user_div,priority:1;index:idx_user_floor,priority:1"`

	// admin user_id who made this punish
	MadeBy int `json:"made_by"`

	// punished because of this floor
	FloorID int `json:"floor_id" gorm:"index:idx_user_floor,priority:2"`

	Floor *Floor `json:"floor"` // foreign key

	DivisionID int `json:"division_id" gorm:"index:idx_user_div,priority:2"`

	Division *Division `json:"division"` // foreign key

	// reason
	Reason string `json:"reason" gorm:"size:128"`
}

type Punishments []*Punishment
