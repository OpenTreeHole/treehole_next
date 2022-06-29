package schemas

import "treehole_next/models"

type DeleteDivision struct {
	// Admin only
	// ID of the target division that all the deleted division's holes will be moved to
	To int `json:"to" default:"1"`
}

type AddDivision struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModifyDivision struct {
	AddDivision
	Pinned models.IntArray `json:"pinned"`
}
