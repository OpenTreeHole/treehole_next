package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"gorm.io/gorm"
)

type Report struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created"`
	UpdatedAt time.Time `json:"time_updated"`
	ReportID  int       `json:"report_id" gorm:"-:all"`
	FloorID   int       `json:"floor_id"`
	HoleID    int       `json:"hole_id" gorm:"-:all"`
	Floor     *Floor    `json:"floor" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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

func (report *Report) Preprocess(c *fiber.Ctx) (err error) {
	err = report.Floor.SetDefaults(c)
	if err != nil {
		return err
	}
	for i := range report.Floor.Mention {
		err = report.Floor.Mention[i].SetDefaults(c)
		if err != nil {
			return err
		}
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
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	existingReport := Report{}
	err = tx.Where("user_id = ? AND floor_id = ?", userID, report.FloorID).First(&existingReport).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		report.UserID = userID
		err = tx.Create(&report).Error

		report.ReportID = report.ID

		err = tx.Model(report).Association("Floor").Find(&report.Floor)
		if err != nil {
			return err
		}

		err = report.Preprocess(c)
		if err != nil {
			return err
		}
	} else {
		existingReport.Reason = existingReport.Reason + "\n" + report.Reason
		err = tx.Model(&existingReport).Update("reason", existingReport.Reason).Error // update reason and load floor in AfterUpdate hook
		report.Floor = existingReport.Floor
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

	//err = report.Preprocess(nil)
	//if err != nil {
	//	return err
	//}

	return nil
}

//var adminCounter = new(int32)

func (report *Report) SendModify(_ *gorm.DB) error {
	// get recipients
	userIDs := []int{report.UserID}

	// construct message
	message := Notification{
		Data:       report,
		Recipients: userIDs,
		Description: fmt.Sprintf(
			"处理结果：%s\n感谢您为维护社区秩序所做的贡献。",
			report.Result,
		),
		Title: "您的举报已得到处理",
		Type:  MessageTypeReportDealt,
		URL:   fmt.Sprintf("/api/reports/%d", report.ID),
	}

	// send
	_, err := message.Send()
	return err
}
