package report

import (
	"fmt"
	"gorm.io/gorm"
	. "treehole_next/models"
)

type Range int

const (
	RangeNotDealt Range = iota
	RangeDealt
	RangeAll
)

type ListModel struct {
	Query
	// Sort order, default is desc
	Sort string `query:"sort" default:"desc" validate:"oneof=asc desc"`
	// Range, 0: not dealt, 1: dealt, 2: all
	Range Range `json:"range"`
}

func (q *ListModel) BaseQuery() *gorm.DB {
	return DB.
		Limit(q.Size).
		Offset(q.Offset).
		Order(fmt.Sprintf("`report`.`%s` %s", q.OrderBy, q.Sort))
}

type AddModel struct {
	FloorID int    `json:"floor_id" validate:"required"`
	Reason  string `json:"reason" validate:"required,max=128"`
}

type DeleteModel struct {
	// The deal result, send it to reporter
	Result string `json:"result" validate:"required,max=128"`
}
