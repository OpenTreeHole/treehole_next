package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type BaseModel struct {
	ID        int       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type IntArray []int

func (p IntArray) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *IntArray) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type StringMap map[string]interface{}

func (p StringMap) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *StringMap) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type IntStringMap map[int]string

func (p IntStringMap) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *IntStringMap) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}
