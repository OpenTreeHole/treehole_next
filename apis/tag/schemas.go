package tag

type CreateModel struct {
	Name string `json:"name,omitempty" validate:"max=10"` // Admin only
}

type ModifyModel struct {
	CreateModel
	Temperature int `json:"temperature,omitempty"` // Admin only
}

type DeleteModel struct {
	// Admin only
	// Name of the target tag that all the deleted tag's holes will be connected to
	To string `json:"to,omitempty"`
}

type SearchModel struct {
	Search string `json:"s" query:"s" validate:"max=32"` // search tag by name
}
