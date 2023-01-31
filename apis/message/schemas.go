package message

type CreateModel struct {
	// MessageTypeMail
	Description string `json:"description"`
	Recipients  []int  `json:"recipients" validate:"required"`
}

type ListModel struct {
	NotRead bool `default:"false" query:"not_read"`
}
