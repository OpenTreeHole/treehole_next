package report

import (
	"github.com/gofiber/fiber/v2"
	. "treehole_next/models"
)

var _ Report

// GetReport
// @Summary Get A Report
// @Tags Report
// @Produce application/json
// @Router /reports/{id} [get]
// @Param id path int true "id"
// @Success 200 {object} Report
// @Failure 404 {object} MessageModel
func GetReport(c *fiber.Ctx) error {
	return nil
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
	return nil
}

// AddReport
// @Summary Add a report
// @Description Add a report and send notification to admins
// @Tags Report
// @Produce application/json
// @Router /reports [post]
// @Param json body AddModel true "json"
// @Success 200 {object} Report
// @Failure 400 {object} utils.HttpError
func AddReport(c *fiber.Ctx) error {
	return nil
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
	return nil
}
