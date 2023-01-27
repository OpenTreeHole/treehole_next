package models

import (
	"time"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Division struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"time_created"`
	UpdatedAt   time.Time `json:"time_updated"`
	DivisionID  int       `json:"division_id" gorm:"-:all"`
	Name        string    `json:"name" gorm:"unique"`
	Description string    `json:"description"`
	Pinned      IntArray  `json:"-"` // pinned holes in given order
	Holes       []Hole    `json:"pinned"`
}

func (division Division) GetID() int {
	return division.ID
}

type Divisions []*Division

func (divisions Divisions) Preprocess(c *fiber.Ctx) error {
	for i := 0; i < len(divisions); i++ {
		err := divisions[i].Preprocess(c)
		if err != nil {
			return err
		}
	}
	return utils.SetCache("divisions", divisions, 0)
}

func (division *Division) Preprocess(c *fiber.Ctx) error {
	var pinned = []int(division.Pinned)
	if len(pinned) == 0 {
		division.Holes = []Hole{}
		return nil
	}
	var holes Holes
	DB.Find(&holes, pinned)
	holes = utils.OrderInGivenOrder(holes, pinned)
	err := holes.Preprocess(c)
	if err != nil {
		return err
	}
	division.Holes = holes
	return nil
}

func (division *Division) AfterFind(tx *gorm.DB) (err error) {
	division.DivisionID = division.ID
	return nil
}

func (division *Division) AfterCreate(tx *gorm.DB) (err error) {
	division.DivisionID = division.ID
	return nil
}
