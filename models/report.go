package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Report struct {
	BaseModel
	FloorID int    `json:"floor_id"`
	Floor   Floor  `json:"floor"`
	UserID  int    `json:"-"` // the reporter's id, should keep a secret
	Reason  string `json:"reason" gorm:"size:128"`
	Dealt   bool   `json:"dealt"`                  // the report has been dealt
	DealtBy int    `json:"dealt_by"`               // who dealt the report
	Result  string `json:"result" gorm:"size:128"` // deal result
}

func (report *Report) Create(c *fiber.Ctx, db ...*gorm.DB) error {
	var tx *gorm.DB
	if len(db) > 0 {
		tx = db[0]
	} else {
		tx = DB
	}

	userID, err := GetUserID(c)
	if err != nil {
		return err
	}
	report.UserID = userID
	tx.Create(&report)
	return nil
}
