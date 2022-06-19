package models

import (
	"gorm.io/gorm"
	"treehole_next/config"
)

type HoleFloor struct {
	FirstFloor *Floor   `json:"first_floor"`
	LastFloor  *Floor   `json:"last_floor"`
	Floors     *[]Floor `json:"floors"`
}

type Hole struct {
	BaseModel
	DivisionID int          `json:"division_id"`
	Tags       []*Tag       `json:"tags" gorm:"many2many:hole_tags"`
	Floors     []Floor      `json:"-"`
	HoleFloor  HoleFloor    `json:"floors" gorm:"-:all"` // return floors
	View       int          `json:"view"`
	Reply      int          `json:"reply"`
	Mapping    IntStringMap `json:"-"`
	Hidden     bool         `json:"hidden"`
}
type Holes []Hole

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

func (hole *Hole) Preprocess() error {
	var floors []Floor
	result := DB.Where("hole_id = ?", hole.ID).Limit(config.Config.Size).Find(&floors)
	hole.HoleFloor.Floors = &floors
	if result.RowsAffected > 0 {
		hole.HoleFloor.FirstFloor = &floors[0]
	}

	var floor Floor
	result = DB.Where("hole_id = ?", hole.ID).Last(&floor)
	if result.Error == nil { // last floor exists
		hole.HoleFloor.LastFloor = &floor
	}

	return nil
}
func (holes Holes) Preprocess() error {
	for i := 0; i < len(holes); i++ {
		err := holes[i].Preprocess()
		if err != nil {
			return err
		}
	}
	return nil
}
