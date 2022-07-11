package floor

type content struct {
	// Owner or admin, the original content should be moved to  floor_history
	Content string `json:"content" validate:"required"`
}

type ListOldModel struct {
	HoleID int `query:"hole_id"`
	Size   int `query:"length" default:"10" validate:"min=0,max=30"`
	Offset int `query:"start_floor"`
}

type CreateModel struct {
	content
	// id of the floor to which replied
	ReplyTo int `json:"reply_to" validate:"min=0"`
}

type CreateOldModel struct {
	HoleID int `json:"hole_id" validate:"min=1"`
	CreateModel
}

type ModifyModel struct {
	content
	// All user, deprecated, "add" is like, "cancel" is reset
	Like string `json:"like" validate:"oneof=add cancel"`
	// Admin only
	Fold string `json:"fold" validate:"max=16"`
	// Admin only
	SpecialTag string `json:"special_tag" validate:"max=16"`
}

type DeleteModel struct {
	Reason string `json:"delete_reason" validate:"max=32"`
}
