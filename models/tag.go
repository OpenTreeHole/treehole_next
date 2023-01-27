package models

import (
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"time_created"`
	UpdatedAt   time.Time `json:"time_updated"`
	TagID       int       `json:"tag_id" gorm:"-:all"`
	Name        string    `json:"name,omitempty" gorm:"unique;size:32"`
	Temperature int       `json:"temperature,omitempty"`
	Holes       []*Hole   `json:"-" gorm:"many2many:hole_tags"`
}

func (tag Tag) GetID() int {
	return tag.ID
}

func (tag *Tag) AfterFind(tx *gorm.DB) (err error) {
	tag.TagID = tag.ID
	return nil
}

func (tag *Tag) AfterCreate(tx *gorm.DB) (err error) {
	tag.TagID = tag.ID
	return nil
}
