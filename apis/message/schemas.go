package message

import (
	. "treehole_next/models"

	"github.com/creasty/defaults"
)

type CreateModel struct {
	// message type, change "oneof" when MessageType changes
	Type        MessageType `json:"type" validate:"required,oneof=favorite reply mention modify report permission report_dealt"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Data        JSON        `json:"data"`
	URL         string      `json:"url"`
	Recipients  []int       `json:"recipients" validate:"required"`
}

type ListModel struct {
	NotRead bool `default:"false" query:"not_read"`
}

func (body *CreateModel) SetDefaults() {
	if defaults.CanUpdate(body.Data) {
		body.Data = JSON{}
	}
}
