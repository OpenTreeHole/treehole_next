package schemas

type CreateTag struct {
	Name string `json:"name,omitempty" gorm:"unique;size:32"` // Admin only
}

type ModifyTag struct {
	CreateTag
	Temperature int `json:"temperature,omitempty"` // Admin only
}

type DeleteTag struct {
	// Admin only
	// Name of the target tag that all the deleted tag's holes will be connected to
	To string `json:"to,omitempty"`
}
