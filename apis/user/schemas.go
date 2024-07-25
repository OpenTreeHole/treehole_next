package user

type ModifyModel struct {
	Nickname *string          `json:"nickname" validate:"omitempty,min=1"`
	Config   *UserConfigModel `json:"config"`
}

type UserConfigModel struct {
	Notify     []string `json:"notify"`
	ShowFolded *string  `json:"show_folded"`
}
