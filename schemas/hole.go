package schemas

import "time"

type GetHoleOld struct {
	Offset     time.Time `query:"start_time"`
	Size       int       `query:"length,omitempty"`
	DivisionID int       `query:"division_id,omitempty"`
	Tag        string    `query:"tag,omitempty"`
}

type tags struct {
	Tags []CreateTag // All users
}

type divisionID struct {
	DivisionID int `json:"division_id"` // Admin only
}

type CreateHole struct {
	Content string `json:"content"`
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
