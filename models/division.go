package models

import (
	"github.com/gofiber/fiber/v2"
	"treehole_next/utils"
)

type Division struct {
	BaseModel
	Name        string   `json:"name" gorm:"unique" `
	Description string   `json:"description"`
	Pinned      IntArray `json:"-"     ` // pinned holes in given order
	Holes       []Hole   `json:"pinned"     `
}

type Divisions []*Division

func (divisions Divisions) Preprocess(c *fiber.Ctx) error {
	for _, d := range divisions {
		err := d.Preprocess(c)
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
