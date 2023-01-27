package models

import (
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	/// saved fields
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"-" gorm:"not null"`
	UpdatedAt time.Time `json:"-" gorm:"not null"`

	/// base info
	Name        string `json:"name" gorm:"not null;unique;size:32"`
	Temperature int    `json:"temperature" gorm:"not null;default:0"`

	/// association info, should add foreign key
	Holes []*Hole `json:"-" gorm:"many2many:hole_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	/// generated field
	TagID int `json:"tag_id" gorm:"-:all"`
}

func (tag Tag) GetID() int {
	return tag.ID
}

func (tag *Tag) AfterFind(tx *gorm.DB) (err error) {
	_ = tx
	tag.TagID = tag.ID
	return nil
}

func (tag *Tag) AfterCreate(tx *gorm.DB) (err error) {
	_ = tx
	tag.TagID = tag.ID
	return nil
}
