// Package schemas contains structs for validation and documentation
package schemas

import "time"

type Query struct {
	Size    int    `query:"size" default:"10"`     // length of object array
	Offset  int    `query:"offset" default:"0"`    // offset of object array
	OrderBy string `query:"order_by" default:"id"` // Now only supports id
	Desc    bool   `query:"desc" default:"false"`  // Is descending order
}

type QueryTime struct {
	Size int `json:"size" default:"10"` // length of object array
	// updated time < offset (default is now)
	Offset time.Time `json:"offset"`
}

type MessageModel struct {
	Message string `json:"message,omitempty"`
}
