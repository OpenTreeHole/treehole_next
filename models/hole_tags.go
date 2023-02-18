package models

type HoleTag struct {
	HoleID int `json:"hole_id"`
	TagID  int `json:"tag_id"`
}

func (HoleTag) TableName() string {
	return "hole_tags"
}

type HoleTags []*HoleTag
