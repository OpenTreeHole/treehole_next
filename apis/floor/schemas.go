package floor

import (
	"treehole_next/models"

	"gorm.io/gorm"
)

type ListModel struct {
	models.Query
	OrderBy string `json:"order_by" query:"order_by" default:"id" validate:"oneof=id like"` // SQL ORDER BY field
}

func (q *ListModel) BaseQuery() *gorm.DB {
	return models.DB.Limit(q.Size).Offset(q.Offset).Order(q.OrderBy + " " + q.Sort)
}

type ListOldModel struct {
	HoleID int     `query:"hole_id"     json:"hole_id"`
	Size   int     `query:"length"      json:"length"     validate:"min=0,max=50" `
	Offset int     `query:"start_floor" json:"start_floor"`
	Search *string `query:"s"           json:"s"`
}

func (q *ListOldModel) BaseQuery() *gorm.DB {
	return models.DB.Limit(q.Size).Offset(q.Offset)
}

type CreateModel struct {
	Content string `json:"content" validate:"required"`
	// Admin and Operator only
	SpecialTag string `json:"special_tag" validate:"omitempty,max=16"`
	// id of the floor to which replied
	ReplyTo int `json:"reply_to" validate:"min=0"`
}

type CreateOldModel struct {
	HoleID int `json:"hole_id" validate:"min=1"`
	CreateModel
}

type CreateOldResponse struct {
	Data    models.Floor `json:"data"`
	Message string       `json:"message"`
}

type ModifyModel struct {
	// Owner or admin, the original content should be moved to  floor_history
	Content string `json:"content" validate:"omitempty"`
	// Admin and Operator only
	SpecialTag string `json:"special_tag" validate:"omitempty,max=16"`
	// All user, deprecated, "add" is like, "cancel" is reset
	Like string `json:"like" validate:"omitempty,oneof=add cancel"`
	// Admin and operator only, only string, for version 2
	Fold string `json:"fold_v2" validate:"omitempty,max=64"`
	// Admin and operator only, string array, for version 1: danxi app
	FoldFrontend []string `json:"fold" validate:"omitempty"`
}

type DeleteModel struct {
	Reason string `json:"delete_reason" validate:"max=32"`
}

type RestoreModel struct {
	Reason string `json:"restore_reason" validate:"required,max=32"`
}

type SearchConfigModel struct {
	Open bool `json:"open"`
}
