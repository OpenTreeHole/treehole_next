package models

import (
	"gorm.io/gorm"
)

type CanQuery interface {
	BaseQuery() *gorm.DB
}

type Query struct {
	Size    int    `query:"size" default:"30" validate:"min=0,max=50"`    // length of object array
	Offset  int    `query:"offset" default:"0" validate:"min=0"`          // offset of object array
	Sort    string `query:"sort" default:"asc" validate:"oneof=asc desc"` // Sort order
	OrderBy string `query:"order_by" default:"id"`                        // SQL ORDER BY field
}

func (q *Query) BaseQuery() *gorm.DB {
	return DB.Limit(q.Size).Offset(q.Offset).Order(q.OrderBy + " " + q.Sort)
}

type MessageModel struct {
	Message string `json:"message"`
}
