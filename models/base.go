// Package models contains database models
package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Map = map[string]interface{}

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
