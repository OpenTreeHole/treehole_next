// Package models contains database models
package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Map = map[string]interface{}

type BaseModel struct {
	ID        int       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"time_created"`
	UpdatedAt time.Time `json:"time_updated"`
}

type IntArray []int

func (p IntArray) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *IntArray) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type Models interface {
	Division | Hole | Floor | Tag | User | Report |
		[]Division | []Hole | []Floor | []Tag | []User | []Report
}

type Model interface {
	Division | Hole | Floor | Tag | User
}
