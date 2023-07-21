package models

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
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
	DealtBy int    `json:"dealt_by" gorm:"index"`
	Result  string `json:"result" gorm:"size:128"` // deal result
}

func (report *Report) GetID() int {
	return report.ID
}

type Reports []*Report

func (report *Report) Preprocess(_ *fiber.Ctx) error {
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

	err = tx.Model(report).Association("Floor").Find(&report.Floor)
	if err != nil {
		return err
	}

	err = report.Preprocess(nil)
	if err != nil {
		return err
	}

	return nil
}

func (report *Report) AfterFind(_ *gorm.DB) (err error) {
	report.ReportID = report.ID

	return nil
}

func (report *Report) AfterUpdate(tx *gorm.DB) (err error) {
	err = tx.Model(report).Association("Floor").Find(&report.Floor)
	if err != nil {
		return err
	}

	err = report.Preprocess(nil)
	if err != nil {
		return err
	}

	return nil
}

var adminCounter = new(int32)

func (report *Report) SendCreate(_ *gorm.DB) error {
	adminList.RLock()
	defer adminList.RUnlock()
	if len(adminList.data) == 0 {
		return nil
	}

	// get counter
	currentCounter := atomic.AddInt32(adminCounter, 1)
	result := atomic.CompareAndSwapInt32(adminCounter, int32(len(adminList.data)), 0)
	if result {
		log.Info().Str("model", "get admin").Msg("adminCounter Reset")
	}
	userIDs := []int{adminList.data[currentCounter-1]}

	// construct message
	message := Notification{
		Data:       report,
		Recipients: userIDs,
		Description: fmt.Sprintf(
			"理由：%s，内容：%s",
			report.Reason,
			report.Floor.Content,
		),
		Title: "您有举报需要处理",
		Type:  MessageTypeReport,
		URL:   fmt.Sprintf("/api/reports/%d", report.ID),
	}

	// send
	_, err := message.Send()
	return err
}

func (report *Report) SendModify(_ *gorm.DB) error {
	// get recipients
	userIDs := []int{report.UserID}

	// construct message
	message := Notification{
		Data:       report,
		Recipients: userIDs,
		Description: fmt.Sprintf(
			"结果：%s，内容：%s",
			report.Result,
			report.Floor.Content,
		),
		Title: "您的举报被处理了",
		Type:  MessageTypeReportDealt,
		URL:   fmt.Sprintf("/api/reports/%d", report.ID),
	}

	// send
	_, err := message.Send()
	return err
}
