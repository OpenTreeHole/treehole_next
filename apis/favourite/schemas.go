package favourite

import . "treehole_next/models"

type Response struct {
	Message string   `json:"message"`
	Data    IntArray `json:"data"`
}

type AddModel struct {
	HoleID int `json:"hole_id"`
}

type ModifyModel struct {
	HoleIDs []int `json:"hole_ids"`
}

type DeleteModel struct {
	HoleID int `json:"hole_id"`
}
