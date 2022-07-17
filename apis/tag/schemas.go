package tag

type CreateModel struct {
	Name string `json:"name,omitempty" gorm:"unique;size:32"` // Admin only
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
