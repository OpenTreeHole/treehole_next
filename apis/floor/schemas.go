package floor

import (
	"time"

	"github.com/opentreehole/go-common"

	"treehole_next/models"
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
	// 仅管理员，留空则重置，高优先级
	Fold *string `json:"fold_v2" validate:"omitempty,max=64"`
	// 仅管理员，留空则重置，低优先级
	FoldFrontend []string `json:"fold" validate:"omitempty"`
}

func (body ModifyModel) DoNothing() bool {
	return body.Content == nil && body.SpecialTag == nil && body.Like == nil && body.Fold == nil && body.FoldFrontend == nil
}

func (body ModifyModel) CheckPermission(user *models.User, floor *models.Floor, hole *models.Hole) error {
	if body.Content != nil {
		if !user.IsAdmin {
			if user.ID != floor.UserID {
				return common.Forbidden("这不是您的楼层，您没有权限修改")
			} else {
				if user.BanDivision[hole.DivisionID] != nil {
					return common.Forbidden(user.BanDivisionMessage(hole.DivisionID))
				} else if hole.Locked {
					return common.Forbidden("此洞已被锁定，您无法修改")
				} else if floor.Deleted {
					return common.Forbidden("此洞已被删除，您无法修改")
				}
			}
		} else {
			if user.BanDivision[hole.DivisionID] != nil {
				return common.Forbidden(user.BanDivisionMessage(hole.DivisionID))
			}
		}
	}
	if (body.Fold != nil || body.FoldFrontend != nil) && !user.IsAdmin {
		return common.Forbidden("非管理员禁止折叠")
	}
	if body.SpecialTag != nil && !user.IsAdmin {
		return common.Forbidden("非管理员禁止修改特殊标签")
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

type SensitiveFloorRequest struct {
	Size    int               `json:"size" query:"size" default:"10" validate:"max=10"`
	Offset  common.CustomTime `json:"offset" query:"offset" swaggertype:"string"`
	OrderBy string            `json:"order_by" query:"order_by" default:"time_created" validate:"oneof=time_created time_updated"`
	Open    bool              `json:"open" query:"open"`
	All     bool              `json:"all" query:"all"`
}

type SensitiveFloorResponse struct {
	ID                int       `json:"id"`
	CreatedAt         time.Time `json:"time_created"`
	UpdatedAt         time.Time `json:"time_updated"`
	Content           string    `json:"content"`
	Modified          int       `json:"modified"`
	IsActualSensitive *bool     `json:"is_actual_sensitive"`
	HoleID            int       `json:"hole_id"`
	Deleted           bool      `json:"deleted"`
	SensitiveDetail   string    `json:"sensitive_detail,omitempty"`
}

func (s *SensitiveFloorResponse) FromModel(floor *models.Floor) *SensitiveFloorResponse {
	s.ID = floor.ID
	s.CreatedAt = floor.CreatedAt
	s.UpdatedAt = floor.UpdatedAt
	s.Content = floor.Content
	s.Modified = floor.Modified
	s.IsActualSensitive = floor.IsActualSensitive
	s.HoleID = floor.HoleID
	s.Deleted = floor.Deleted
	s.SensitiveDetail = floor.SensitiveDetail
	return s
}

type ModifySensitiveFloorRequest struct {
	IsActualSensitive bool `json:"is_actual_sensitive"`
}

type BanDivision map[int]*time.Time
