package models

type Report struct {
	BaseModel
	FloorID int    `json:"floor_id"`
	Floor   Floor  `json:"floor"`
	UserID  int    `json:"-"` // the reporter's id, should keep a secret
	Reason  string `json:"reason" gorm:"size:128"`
	Dealt   bool   `json:"dealt"` // the report has been dealt
}
