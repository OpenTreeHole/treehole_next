package hole

import (
	"github.com/opentreehole/go-common"
	"time"
	"treehole_next/apis/tag"
	"treehole_next/models"
)

type QueryTime struct {
	Size int `json:"size" query:"size" default:"10" validate:"max=10"`
	// updated time < offset (default is now)
	Offset common.CustomTime `json:"offset" query:"offset" swaggertype:"string"`
	Order  string            `json:"order" query:"order"`
}

func (q *QueryTime) SetDefaults() {
	if q.Offset.IsZero() {
		q.Offset = common.CustomTime{Time: time.Now()}
	}
}

type ListOldModel struct {
	Offset     common.CustomTime `json:"start_time" query:"start_time" swaggertype:"string"`
	Size       int               `json:"length" query:"length" default:"10" validate:"max=10" `
	Tag        string            `json:"tag" query:"tag"`
	DivisionID int               `json:"division_id" query:"division_id"`
	Order      string            `json:"order" query:"order"`
}

func (q *ListOldModel) SetDefaults() {
	if q.Offset.IsZero() {
		q.Offset = common.CustomTime{Time: time.Now()}
	}
}

type TagCreateModelSlice struct {
	Tags []tag.CreateModel `json:"tags" validate:"omitempty,min=1,max=10,dive"` // All users
}

func (tagCreateModelSlice TagCreateModelSlice) ToName() []string {
	tags := make([]string, 0, len(tagCreateModelSlice.Tags))
	for _, tagCreateModel := range tagCreateModelSlice.Tags {
		tags = append(tags, tagCreateModel.Name)
	}
	return tags
}

type CreateModel struct {
	Content string `json:"content" validate:"required"`
	TagCreateModelSlice
	// Admin and Operator only
	SpecialTag string `json:"special_tag" validate:"max=16"`
}

type CreateOldModel struct {
	CreateModel
	DivisionID int `json:"division_id" validate:"omitempty,min=1" default:"1"`
}

type CreateOldResponse struct {
	Data    models.Hole `json:"data"`
	Message string      `json:"message"`
}

type ModifyModel struct {
	TagCreateModelSlice
	DivisionID *int  `json:"division_id" validate:"omitempty,min=1"` // Admin and owner only
	Unhidden   *bool `json:"unhidden"`                               // admin only
	Lock       *bool `json:"lock"`                                   // admin only
}

func (body ModifyModel) CheckPermission(user *models.User, hole *models.Hole) error {
	if body.DivisionID != nil && !user.IsAdmin {
		return common.Forbidden("非管理员禁止修改分区")
	}
	if body.Unhidden != nil && !user.IsAdmin {
		return common.BadRequest("非管理员禁止取消隐藏")
	}
	if body.Tags != nil && !(user.IsAdmin || user.ID == hole.UserID) {
		return common.Forbidden()
	}
	if body.Tags != nil && len(body.Tags) == 0 {
		return common.BadRequest("tags 不能为空")
	}
	if body.Lock != nil && !user.IsAdmin {
		return common.Forbidden("非管理员禁止锁定帖子")
	}
	return nil
}

func (body ModifyModel) DoNothing() bool {
	return body.Unhidden == nil && body.Tags == nil && body.DivisionID == nil && body.Lock == nil
}
