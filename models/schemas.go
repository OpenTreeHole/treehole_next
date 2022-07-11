package models

import "gorm.io/gorm"

type Query struct {
	Size    int    `query:"size" default:"10" validate:"min=0,max=30"`                 // length of object array
	Offset  int    `query:"offset" default:"0" validate:"min=0"`                       // offset of object array
	OrderBy string `query:"order_by" default:"storey" validate:"oneof=storey id like"` // SQL ORDER BY field
	Sort    string `query:"sort" default:"asc" validate:"oneof=asc desc"`              // Sort order
}

type MessageModel struct {
	Message string `json:"message,omitempty"`
}

func BaseQuery(q *Query) *gorm.DB {
	return DB.Limit(q.Size).Offset(q.Offset).Order(q.OrderBy + " " + q.Sort)
}
