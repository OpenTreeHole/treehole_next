package subscription

type Response struct {
	Message string `json:"message"`
	Data    []int  `json:"data"`
}

type ListModel struct {
	Plain bool `json:"plain" default:"false" query:"plain"`
}

type AddModel struct {
	HoleID int `json:"hole_id"`
}

type DeleteModel struct {
	HoleID int `json:"hole_id"`
}
