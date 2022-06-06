package models

import (
	"gorm.io/gorm"
)

type Hole struct {
	BaseModel
	DivisionID int          `json:"division_id,omitempty"`
	Tags       []*Tag       `json:"tags,omitempty" gorm:"many2many:hole_tags"`
	Floors     []Floor      `json:"floors,omitempty"`
	View       int          `json:"view,omitempty"`
	Reply      int          `json:"reply,omitempty"`
	Mapping    IntStringMap `json:"mapping,omitempty"`
	Hidden     bool         `json:"hidden,omitempty"`
}

// AfterFind set default mapping as {}
//goland:noinspection GoUnusedParameter
func (hole *Hole) AfterFind(tx *gorm.DB) (err error) {
	if hole.Mapping == nil {
		hole.Mapping = map[int]string{}
	}
	return
}

// AfterCreate set default mapping as {}
func (hole *Hole) AfterCreate(tx *gorm.DB) (err error) {
	return hole.AfterFind(tx)
}
