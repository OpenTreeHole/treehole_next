package favourite

type AddModel struct {
	HoleID int `json:"hole_id"`
}

type ModifyModel struct {
	HoleIDs []int `json:"hole_ids"`
}

type DeleteModel struct {
	HoleID int `json:"hole_id"`
}
