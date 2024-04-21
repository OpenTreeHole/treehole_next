package report

import (
	"fmt"
	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"
	"time"
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetReport
//
// @Summary Get A Report
// @Tags Report
// @Produce application/json
// @Router /reports/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Report
// @Failure 404 {object} MessageModel
func GetReport(c *fiber.Ctx) error {
	// validate query
	reportID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// find report
	var report Report
	result := LoadReportFloor(DB).First(&report, reportID)
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &report)
}

// ListReports
//
// @Summary List All Reports
// @Tags Report
// @Produce application/json
// @Router /reports [get]
// @Param object query ListModel false "query"
// @Success 200 {array} Report
// @Failure 404 {object} MessageModel
func ListReports(c *fiber.Ctx) error {
	// validate query
	var query ListModel
	err := common.ValidateQuery(c, &query)
	if err != nil {
		return err
	}

	// find reports
	var reports Reports

	querySet := LoadReportFloor(query.BaseQuery())

	var result *gorm.DB
	switch query.Range {
	case RangeNotDealt:
		result = querySet.Find(&reports, "dealt = ?", false)
	case RangeDealt:
		result = querySet.Find(&reports, "dealt = ?", true)
	case RangeAll:
		result = querySet.Find(&reports)
	}
	if result.Error != nil {
		return result.Error
	}
	return Serialize(c, &reports)
}

// AddReport
//
// @Summary Add a report
// @Description Add a report and send notification to admins
// @Tags Report
// @Produce application/json
// @Router /reports [post]
// @Param json body AddModel true "json"
// @Success 204
//
// @Failure 400 {object} common.HttpError
func AddReport(c *fiber.Ctx) error {
	// validate body
	var body AddModel
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	user, err := GetUser(c)
	if err != nil {
		if err != nil {
			return err
		}
	}

	// permission
	if user.BanReport != nil {
		return common.Forbidden(user.BanReportMessage())
	}

	// add report
	report := Report{
		FloorID: body.FloorID,
		Reason:  body.Reason,
		Dealt:   false,
	}
	err = report.Create(c)
	if err != nil {
		return err
	}

	// Send Notification
	err = report.SendCreate(DB)
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("SendCreate failed: ")
		// return err // only for test
	}

	return c.Status(204).JSON(nil)
}

// DeleteReport
//
// @Summary Deal a report
// @Description Mark a report as "dealt" and send notification to reporter
// @Tags Report
// @Produce application/json
// @Router /reports/{id} [delete]
// @Param id path int true "id"
// @Param json body DeleteModel true "json"
// @Success 200 {object} Report
// @Failure 400 {object} common.HttpError
func DeleteReport(c *fiber.Ctx) error {
	// validate query
	reportID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// validate body
	var body DeleteModel
	err = common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get user id
	userID, err := common.GetUserID(c)
	if err != nil {
		return err
	}

	// modify report
	var report Report
	result := LoadReportFloor(DB).First(&report, reportID)
	if result.Error != nil {
		return result.Error
	}
	report.Dealt = true
	report.DealtBy = userID
	report.Result = body.Result
	DB.Omit("Floor").Save(&report)

	MyLog("Report", "Delete", reportID, userID, RoleAdmin)

	// Send Notification
	err = report.SendModify(DB)
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("SendModify failed")
		// return err // only for test
	}

	return Serialize(c, &report)
}

type banBody struct {
	Days   *int   `json:"days" validate:"omitempty,min=1"`
	Reason string `json:"reason"` // optional
}

// BanReporter
//
// @Summary Ban reporter of a report
// @Tags Report
// @Produce json
// @Router /reports/ban/{id} [post]
// @Param json body banBody true "json"
// @Success 201 {object} User
func BanReporter(c *fiber.Ctx) error {
	// validate body
	var body banBody
	err := common.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	reportID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// get user
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	// permission
	if !user.IsAdmin {
		return common.Forbidden()
	}

	var report Report
	err = DB.Take(&report, reportID).Error
	if err != nil {
		return err
	}

	var days int
	if body.Days != nil {
		days = *body.Days
		if days <= 0 {
			days = 1
		}
	} else {
		days = 1
	}

	duration := time.Duration(days) * 24 * time.Hour

	reportPunishment := ReportPunishment{
		UserID:   report.UserID,
		MadeBy:   user.ID,
		ReportId: report.ID,
		Duration: &duration,
		Reason:   body.Reason,
	}
	user, err = reportPunishment.Create()
	if err != nil {
		return err
	}

	// construct message for user
	message := Notification{
		Data:       report,
		Recipients: []int{report.UserID},
		Description: fmt.Sprintf(
			"您因违反社区公约被禁止举报。时间：%d天，原因：%s\n如有异议，请联系admin@fduhole.com。",
			days,
			body.Reason,
		),
		Title: "处罚通知",
		Type:  MessageTypePermission,
		URL:   fmt.Sprintf("/api/reports/%d", report.ID),
	}

	// send
	_, err = message.Send()
	if err != nil {
		return err
	}

	return c.JSON(user)
}
