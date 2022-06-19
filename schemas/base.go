// Package schemas contains structs for validation and documentation
package schemas

import "time"

type Query struct {
	Size int `json:"size" default:"10"` // length of object array
	// If a time, order by updated time (for created time, ordering by id is better)
	// Otherwise, the int is passed after sql "offset"
	Offset  string `json:"offset" default:"0"`
	OrderBy string `json:"order_by" default:"id"` // Now only supports id
	Desc    bool   `json:"desc" default:"true"`   // Is descending order
}

type QueryTime struct {
	Size int `json:"size" default:"10"` // length of object array
	// updated time < offset (default is now)
	Offset time.Time `json:"offset"`
}

type MessageModel struct {
	Message string `json:"message,omitempty"`
}
