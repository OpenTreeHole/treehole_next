package models

import (
	"gorm.io/gorm"
)

type Division struct {
	BaseModel
	Name        string   `json:"name" gorm:"unique" `
	Description string   `json:"description"`
	Pinned      IntArray `json:"pinned"     ` // pinned holes in given order
}

// AfterFind set default pinned as []
//goland:noinspection GoUnusedParameter
func (division *Division) AfterFind(tx *gorm.DB) (err error) {
	if division.Pinned == nil {
		division.Pinned = []int{}
	}
	return
}

// AfterCreate set default pinned as []
func (division *Division) AfterCreate(tx *gorm.DB) (err error) {
	return division.AfterFind(tx)
}
