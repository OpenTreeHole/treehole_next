package report

import (
	. "treehole_next/models"
	. "treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetReport
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
	result := DB.Joins("Floor").First(&report, reportID)
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&report)
}

// ListReports
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
	err := ValidateQuery(c, &query)
	if err != nil {
		return err
	}
	query.OrderBy = "`report`.`" + query.OrderBy + "`"

	// find reports
	var reports []Report
	BaseQuerySet := query.BaseQuery().Joins("Floor")
	var result *gorm.DB
	switch query.Range {
	case RangeNotDealt:
		result = BaseQuerySet.Find(&reports, "dealt = ?", false)
	case RangeDealt:
		result = BaseQuerySet.Find(&reports, "dealt = ?", true)
	case RangeAll:
		result = BaseQuerySet.Find(&reports)
	}
	if result.Error != nil {
		return result.Error
	}
	return c.JSON(&reports)
}

// AddReport
// @Summary Add a report
// @Description Add a report and send notification to admins
// @Tags Report
// @Produce application/json
// @Router /reports [post]
// @Param json body AddModel true "json"
// @Success 204
// @Failure 400 {object} utils.HttpError
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

	// TODO: notification to admin

	return c.Status(204).JSON(nil)
}

// DeleteReport
// @Summary Deal a report
// @Description Mark a report as "dealt" and send notification to reporter
// @Tags Report
// @Produce application/json
// @Router /reports/{id} [delete]
// @Param id path int true "id"
// @Param json body DeleteModel true "json"
// @Success 200 {object} Report
// @Failure 400 {object} utils.HttpError
func DeleteReport(c *fiber.Ctx) error {
	// validate query
	reportID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// modify report
	var report Report
	result := DB.Joins("Floor").First(&report, reportID)
	if result.Error != nil {
		return result.Error
	}
	report.Dealt = true

	// save report
	DB.Omit("Floor").Save(&report)

	// TODO: notification to reporter

	return c.JSON(&report)
}
