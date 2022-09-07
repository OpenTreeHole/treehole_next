package hole

import (
	"time"
	"treehole_next/apis/tag"
	"treehole_next/models"
)

type QueryTime struct {
	Size int `json:"size" default:"10" validate:"max=10"`
	// updated time < offset (default is now)
	Offset time.Time `json:"offset"`
}

func (q *QueryTime) SetDefaults() {
	if q.Offset.IsZero() {
		q.Offset = time.Now()
	}
}

type ListOldModel struct {
	Offset     time.Time `json:"start_time" query:"start_time"`
	Size       int       `json:"length"     query:"length"      default:"10" validate:"max=10" `
	Tag        string    `json:"tag"        query:"tag"`
	DivisionID int       `json:"division_id" query:"division_id"`
}

func (q *ListOldModel) SetDefaults() {
	if q.Offset.IsZero() {
		q.Offset = time.Now()
	}
}

type tags struct {
	Tags []tag.CreateModel // All users
}

type divisionID struct {
	DivisionID int `json:"division_id" validate:"omitempty,min=1"` // Admin only
}

type CreateModel struct {
	Content string `json:"content" validate:"required"`
	tags
	// Admin and Operator only
	SpecialTag string `json:"special_tag" validate:"max=16"`
}

type CreateOldModel struct {
	CreateModel
	divisionID
}

type CreateOldResponse struct {
	Data    models.Hole
	Message string
}

type ModifyModel struct {
	tags
	divisionID
}
