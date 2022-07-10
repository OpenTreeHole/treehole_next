package floor

type ListOldModel struct {
	HoleID int `query:"hole_id"`
	Size   int `query:"length" default:"20"`
	Offset int `query:"start_floor"  default:"0"`
}

type CreateModel struct {
	Content string `json:"content"`
	// id of the floor to which replied
	ReplyTo int `json:"reply_to"`
}

type CreateOldModel struct {
	HoleID int `json:"hole_id"`
	CreateModel
}

type ModifyModel struct {
	// Owner or admin, the original content should be moved to  floor_history
	Content string `json:"content"`
	// All user, deprecated, "add" is like, "cancel" is reset
	Like string `json:"like"`
	// Admin only
	Fold string `json:"fold"`
	// Admin only
	SpecialTag string `json:"special_tag" default:""`
}

type DeleteModel struct {
	Reason string `json:"delete_reason"`
}
