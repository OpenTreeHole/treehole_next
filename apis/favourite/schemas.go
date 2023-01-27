package favourite

type Response struct {
	Message string `json:"message"`
	Data    []int  `json:"data"`
}

type ListModel struct {
	Plain bool `default:"false" query:"plain"`
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
