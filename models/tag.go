package models

type Tag struct {
	BaseModel
	Name        string  `json:"name,omitempty" gorm:"unique;size:32"`
	Temperature int     `json:"temperature,omitempty"`
	Holes       []*Hole `json:"holes,omitempty" gorm:"many2many:hole_tags"`
}
