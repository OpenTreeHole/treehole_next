package models

type Tag struct {
	BaseModel
	Name        string  `json:"name,omitempty" gorm:"unique;size:32"`
	Temperature int     `json:"temperature,omitempty"`
	Holes       []*Hole `json:"-" gorm:"many2many:hole_tags"`
}
