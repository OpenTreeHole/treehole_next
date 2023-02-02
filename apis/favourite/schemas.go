package favourite

type Response struct {
	Message string `json:"message"`
	Data    []int  `json:"data"`
}

type ListModel struct {
	Order string `json:"order" query:"order" validate:"omitempty,oneof=id time_created hole_time_updated" default:"time_created"`
	Plain bool   `json:"plain" default:"false" query:"plain"`
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
