package user

import "treehole_next/models"

type ModifyModel struct {
	Nickname *string            `json:"nickname" validate:"omitempty,min=1"`
	Config   *models.UserConfig `json:"config"`
}
