package report

import (
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetReport
//
//	@Summary	Get A Report
//	@Tags		Report
//	@Produce	application/json
//	@Router		/reports/{id} [get]
//	@Param		id	path		int	true	"id"
//	@Success	200	{object}	Report
//	@Failure	404	{object}	MessageModel
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
//	@Summary	List All Reports
//	@Tags		Report
//	@Produce	application/json
//	@Router		/reports [get]
//	@Param		object	query		ListModel	false	"query"
//	@Success	200		{array}		Report
//	@Failure	404		{object}	MessageModel
func ListReports(c *fiber.Ctx) error {
	// validate query
	var query ListModel
	err := ValidateQuery(c, &query)
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
//	@Summary		Add a report
//	@Description	Add a report and send notification to admins
//	@Tags			Report
//	@Produce		application/json
//	@Router			/reports [post]
//	@Param			json	body	AddModel	true	"json"
//	@Success		204
//	@Failure		400	{object}	utils.HttpError
func AddReport(c *fiber.Ctx) error {
	// validate body
	var body AddModel
	err := ValidateBody(c, &body)
	if err != nil {
		return err
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

	return c.Status(204).JSON(nil)
}

// DeleteReport
//
//	@Summary		Deal a report
//	@Description	Mark a report as "dealt" and send notification to reporter
//	@Tags			Report
//	@Produce		application/json
//	@Router			/reports/{id} [delete]
//	@Param			id		path		int			true	"id"
//	@Param			json	body		DeleteModel	true	"json"
//	@Success		200		{object}	Report
//	@Failure		400		{object}	utils.HttpError
func DeleteReport(c *fiber.Ctx) error {
	// validate query
	reportID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// validate body
	var body DeleteModel
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// get user id
	userID, err := GetUserID(c)
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

	return Serialize(c, &report)
}
