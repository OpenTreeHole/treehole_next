package schemas

import "time"

type GetHoleOld struct {
	StartTime  time.Time `json:"start_time"`
	Length     int       `json:"length,omitempty"`
	DivisionID int       `json:"division_id,omitempty"`
	Tag        string    `json:"tag,omitempty"`
}

type tags struct {
	Tags []CreateTag // All users
}

type divisionID struct {
	DivisionID int `json:"division_id,omitempty"` // Admin only
}

type CreateHole struct {
	CreateFloor
	tags
}

type CreateHoleOld struct {
	CreateHole
	divisionID
}

type ModifyHole struct {
	tags
	divisionID
}
