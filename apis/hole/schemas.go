package hole

import (
	"time"
	"treehole_next/apis/tag"
)

type QueryTime struct {
	Size int `json:"size" default:"10"` // length of object array
	// updated time < offset (default is now)
	Offset time.Time `json:"offset"`
}

type ListOldModel struct {
	Offset     time.Time `query:"start_time"`
	Size       int       `query:"length,omitempty"`
	DivisionID int       `query:"division_id,omitempty"`
	Tag        string    `query:"tag,omitempty"`
}

type tags struct {
	Tags []tag.CreateModel // All users
}

type divisionID struct {
	DivisionID int `json:"division_id"` // Admin only
}

type CreateModel struct {
	Content string `json:"content"`
	tags
}

type CreateOldModel struct {
	CreateModel
	divisionID
}

type ModifyModel struct {
	tags
	divisionID
}
