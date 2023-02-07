package floor

import (
	"fmt"
	"treehole_next/models"
	"treehole_next/utils"
)

type ListModel struct {
	Size    int    `json:"size" query:"size" default:"30" validate:"min=0,max=50"`          // length of object array
	Offset  int    `json:"offset" query:"offset" default:"0" validate:"min=0"`              // offset of object array
	Sort    string `json:"sort" query:"sort" default:"asc" validate:"oneof=asc desc"`       // Sort order
	OrderBy string `json:"order_by" query:"order_by" default:"id" validate:"oneof=id like"` // SQL ORDER BY field
}

type ListOldModel struct {
	HoleID int    `query:"hole_id"     json:"hole_id"`
	Size   int    `query:"length"      json:"length"     validate:"min=0,max=50" `
	Offset int    `query:"start_floor" json:"start_floor"`
	Search string `query:"s"           json:"s"`
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
	Content *string `json:"content" validate:"omitempty"`
	// Admin and Operator only
	SpecialTag *string `json:"special_tag" validate:"omitempty,max=16"`
	// All user, deprecated, "add" is like, "cancel" is reset
	Like *string `json:"like" validate:"omitempty,oneof=add cancel"`
	// Admin and operator only, only string, for version 2
	Fold *string `json:"fold_v2" validate:"omitempty,max=64"`
	// Admin and operator only, string array, for version 1: danxi app
	FoldFrontend []string `json:"fold" validate:"omitempty"`
}

func (body ModifyModel) DoNothing() bool {
	return body.Content == nil && body.SpecialTag == nil && body.Like == nil && body.Fold == nil && body.FoldFrontend == nil
}

func (body ModifyModel) CheckPermission(user *models.User, floorUserID int, hole *models.Hole) error {
	if user.BanDivision[hole.DivisionID] != nil {
		return utils.Forbidden(fmt.Sprintf("您在此分区已被禁言，解封时间：%s", user.BanDivision[hole.DivisionID]))
	}
	if body.Content != nil && !(user.IsAdmin || (user.ID == floorUserID && !hole.Locked)) {
		return utils.Forbidden("禁止修改此楼")
	}
	if (body.Fold != nil || body.FoldFrontend != nil || body.SpecialTag != nil) && !user.IsAdmin {
		return utils.Forbidden("非管理员禁止修改")
	}
	return nil
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
