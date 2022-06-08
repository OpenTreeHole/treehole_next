// Package schemas contains structs for validation and documentation
package schemas

import "time"

// DocQuery is for swagger docs only, because it doesn't support generic
type DocQuery struct {
	Size int `json:"size" default:"10"` // length of object array
	// Either a time (ISO formatted) or an int
	// If a time, order by updated time (for created time, ordering by id is better)
	// Otherwise, the int is passed after sql "offset"
	Offset  string `json:"offset" default:"0"`
	OrderBy string `json:"order_by" default:"id"` // Now only supports id
	Desc    bool   `json:"desc" default:"true"`   // Is descending order
}

// Query is the base query model
type Query[Offset int | time.Time] struct {
	Size    int    `query:"size"`
	Offset  Offset `query:"offset"`
	OrderBy string `query:"order_by"`
}

type QueryID = Query[int]
type QueryTime = Query[time.Time]

type MessageModel struct {
	Message string `json:"message,omitempty"`
}
