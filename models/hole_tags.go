package models

type HoleTag struct {
	HoleID int `json:"hole_id" gorm:"index"`
	TagID  int `json:"tag_id" gorm:"index"`
}

func (HoleTag) TableName() string {
	return "hole_tags"
}

type HoleTags []*HoleTag
