package report

import "treehole_next/models"

type Range int

const (
	RangeNotDealt Range = iota
	RangeDealt
	RangeAll
)

type ListModel struct {
	models.Query
	Range Range `json:"range"`
}

type AddModel struct {
	FloorID int    `json:"floor_id" validate:"required"`
	Reason  string `json:"reason" validate:"required,max=128"`
}

type DeleteModel struct {
	// The deal result, send it to reporter
	Result string `json:"result" validate:"required,max=128"`
}
