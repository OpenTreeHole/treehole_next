package division

import "treehole_next/models"

type DeleteModel struct {
	// Admin only
	// ID of the target division that all the deleted division's holes will be moved to
	To int `json:"to" default:"1"`
}

type CreateModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModifyModel struct {
	CreateModel
	Pinned models.IntArray `json:"pinned"`
}
