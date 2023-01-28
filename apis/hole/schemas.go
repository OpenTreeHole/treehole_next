package hole

import (
	"time"
	"treehole_next/apis/tag"
	"treehole_next/models"
)

type QueryTime struct {
	Size int `json:"size" default:"10" validate:"max=10"`
	// updated time < offset (default is now)
	Offset models.CustomTime `json:"offset" swaggertype:"string"`
}

func (q *QueryTime) SetDefaults() {
	if q.Offset.IsZero() {
		q.Offset = models.CustomTime{Time: time.Now()}
	}
}

type ListOldModel struct {
	Offset     models.CustomTime `json:"start_time" query:"start_time" swaggertype:"string"`
	Size       int               `json:"length"     query:"length"      default:"10" validate:"max=10" `
	Tag        string            `json:"tag"        query:"tag"`
	DivisionID int               `json:"division_id" query:"division_id"`
	Order      string            `json:"order" query:"order"`
}

func (q *ListOldModel) SetDefaults() {
	if q.Offset.IsZero() {
		q.Offset = models.CustomTime{Time: time.Now()}
	}
}

type TagCreateModelSlice struct {
	Tags []*tag.CreateModel // All users
}

func (tagCreateModelSlice TagCreateModelSlice) ToTags() models.Tags {
	tags := make(models.Tags, 0, len(tagCreateModelSlice.Tags))
	for _, tagCreateModel := range tagCreateModelSlice.Tags {
		tags = append(tags, &models.Tag{Name: tagCreateModel.Name})
	}
	return tags
}

type divisionID struct {
	DivisionID int `json:"division_id" validate:"omitempty,min=1"` // Admin only
}

type CreateModel struct {
	Content string `json:"content" validate:"required"`
	TagCreateModelSlice
	// Admin and Operator only
	SpecialTag string `json:"special_tag" validate:"max=16"`
}

type CreateOldModel struct {
	CreateModel
	divisionID
}

type CreateOldResponse struct {
	Data    models.Hole `json:"data"`
	Message string      `json:"message"`
}

type ModifyModel struct {
	TagCreateModelSlice
	divisionID
	Unhidden bool `json:"unhidden"`
}
