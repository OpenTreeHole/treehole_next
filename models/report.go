package models

import (
	"fmt"
	"sync/atomic"
	"time"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Report struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created"`
	UpdatedAt time.Time `json:"time_updated"`
	ReportID  int       `json:"report_id" gorm:"-:all"`
	FloorID   int       `json:"floor_id"`
	HoleID    int       `json:"hole_id" gorm:"-:all"`
	Floor     *Floor    `json:"floor"`
	UserID    int       `json:"-"` // the reporter's id, should keep a secret
	Reason    string    `json:"reason" gorm:"size:128"`
	Dealt     bool      `json:"dealt"` // the report has been dealt
	// who dealt the report
	DealtBy int    `json:"dealt_by"`
	Result  string `json:"result" gorm:"size:128"` // deal result
}

func (report Report) GetID() int {
	return report.ID
}

type Reports []Report

func (report *Report) Preprocess(c *fiber.Ctx) error {
	report.Floor.SetDefaults()
	for i := range report.Floor.Mention {
		report.Floor.Mention[i].SetDefaults()
	}
	report.HoleID = report.Floor.HoleID
	return nil
}

func (reports Reports) Preprocess(c *fiber.Ctx) error {
	for i := 0; i < len(reports); i++ {
		_ = reports[i].Preprocess(c)
	}
	return nil
}

func LoadReportFloor(tx *gorm.DB) *gorm.DB {
	return tx.Preload("Floor.Mention").Preload("Floor")
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
	report.ReportID = report.ID
	if config.Config.NotificationUrl == "" {
		return nil
	}

	err = tx.Model(report).Association("Floor").Find(&report.Floor)
	if err != nil {
		return err
	}

	err = report.Preprocess(nil)
	if err != nil {
		return err
	}

	err = report.SendCreate(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendCreate failed: " + err.Error())
		// return err // only for test
	}
	return nil
}

func (report *Report) AfterFind(tx *gorm.DB) (err error) {
	report.ReportID = report.ID

	return nil
}

func (report *Report) AfterUpdate(tx *gorm.DB) (err error) {
	if config.Config.NotificationUrl == "" {
		return nil
	}

	err = tx.Model(report).Association("Floor").Find(&report.Floor)
	if err != nil {
		return err
	}

	err = report.Preprocess(nil)
	if err != nil {
		return err
	}

	err = report.SendModify(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendModify failed: " + err.Error())
		// return err // only for test
	}
	return nil
}

var adminCounter = new(int32)

func (report *Report) SendCreate(tx *gorm.DB) error {
	if len(adminList) == 0 {
		return nil
	}

	// get counter
	currentCounter := atomic.AddInt32(adminCounter, 1)
	result := atomic.CompareAndSwapInt32(adminCounter, int32(len(adminList)), 0)
	if result {
		utils.Logger.Info("[getadmin] adminCounter Reset")
	}
	userIDs := []int{adminList[currentCounter-1]}

	// construct message
	message := Message{
		"data":       report,
		"recipients": userIDs,
		"type":       MessageTypeReport,
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
	// get recipients
	userIDs := []int{report.UserID}

	// construct message
	message := Message{
		"data":       report,
		"recipients": userIDs,
		"type":       MessageTypeReportDealt,
		"url":        fmt.Sprintf("/api/reports/%d", report.ID),
	}

	// send
	err := message.Send()
	if err != nil {
		return err
	}

	return nil
}
