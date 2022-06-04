package division

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"treehole_next/db"
)

type AddDivisionModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Division struct {
	db.BaseModel
	Name        string   `json:"name" gorm:"unique" `
	Description string   `json:"description"`
	Pinned      intArray `json:"pinned"     ` // pinned holes in given order
}

type intArray []int

func (p intArray) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *intArray) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
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
