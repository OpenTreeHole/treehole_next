package models

import "gorm.io/gorm"

type Tag struct {
	BaseModel
	TagID       int     `json:"tag_id" gorm:"-:all"`
	Name        string  `json:"name,omitempty" gorm:"unique;size:32"`
	Temperature int     `json:"temperature,omitempty"`
	Holes       []*Hole `json:"-" gorm:"many2many:hole_tags"`
}

func (tag *Tag) AfterFind(tx *gorm.DB) (err error) {
	tag.TagID = tag.ID
	return nil
}
