package schemas

import "treehole_next/models"

type DivisionResponse struct {
	models.BaseModel
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Pinned      []models.Hole `json:"pinned"`
}

type DeleteDivisionModel struct {
	// Admin only
	// ID of the target division that all the deleted division's holes will be moved to
	To int `json:"to" default:"1"`
}

type AddDivisionModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModifyDivisionModel struct {
	AddDivisionModel
	Pinned models.IntArray `json:"pinned"`
}
