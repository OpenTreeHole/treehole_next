package models

import (
	"gorm.io/gorm"
)

type DeleteDivisionModel struct {
	// ID of the target division that move all the deleted division's holes to
	// default to 1
	To int `json:"to"`
}

type AddDivisionModel struct {
	Name        string `json:"name" gorm:"unique" `
	Description string `json:"description"`
}

type ModifyDivisionModel struct {
	AddDivisionModel
	Pinned IntArray `json:"pinned"     `
}

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
