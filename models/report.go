package models

import (
	"fmt"
	"treehole_next/utils"

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

func (report *Report) AfterCreate(tx *gorm.DB) (err error) {
	err = report.SendCreate(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendCreate failed: " + err.Error())
		// return err // only for test
	}
	return nil
}

func (report *Report) AfterUpdate(tx *gorm.DB) (err error) {
	err = report.SendModify(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendModify failed: " + err.Error())
		// return err // only for test
	}
	return nil
}

func (report *Report) SendCreate(tx *gorm.DB) error {
	// get recipents
	userIDs := []int{report.UserID}

	// construct message
	message := Message{
		"data":       report,
		"recipients": userIDs,
		"type":       string(MessageTypeReport),
		"url":        fmt.Sprintf("/api/reports/%d", report.ID),
	}

	// send
	err := message.Send()
	if err != nil {
		return err
	}

	return nil
}

func (report *Report) SendModify(tx *gorm.DB) error {
	// get recipents
	userIDs := []int{report.UserID}

	// construct message
	message := Message{
		"data":       report,
		"recipients": userIDs,
		"type":       string(MessageTypeReportDealt),
		"url":        fmt.Sprintf("/api/reports/%d", report.ID),
	}

	// send
	err := message.Send()
	if err != nil {
		return err
	}

	return nil
}
