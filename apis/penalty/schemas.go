package penalty

type ModifyModel struct {
	Days        *int   `json:"days" validate:"omitempty,min=1"`
	DivisionIDs []int  `json:"division_ids" validate:"omitempty,min=1"`
	Reason      string `json:"reason"`
}

type PostBody struct {
	PenaltyLevel *int   `json:"penalty_level" validate:"omitempty"`   // low priority, deprecated
	Days         *int   `json:"days" validate:"omitempty,min=1"`      // high priority
	Divisions    []int  `json:"divisions" validate:"omitempty,min=1"` // high priority
	Reason       string `json:"reason"`                               // optional
}
