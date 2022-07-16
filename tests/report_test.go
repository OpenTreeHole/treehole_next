package tests

import (
	"log"
	"strconv"
	"testing"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

const (
	REPORT_BASE_ID       = 1
	REPORT_FLOOR_BASE_ID = 1001
)

func init() {
	floors := make([]Floor, 20)
	for i := range floors {
		floors[i].ID = REPORT_FLOOR_BASE_ID + i
	}
	reports := make([]Report, 10)
	for i := range reports {
		reports[i].ID = REPORT_BASE_ID + i
		reports[i].FloorID = REPORT_FLOOR_BASE_ID + i
		if i < 5 {
			reports[i].Dealt = true
		}
	}

	DB.Create(&floors)
	DB.Create(&reports)
}

func TestGetReport(t *testing.T) {
	reportID := REPORT_BASE_ID
	var report Report
	DB.First(&report, reportID)

	var getReport Report
	testAPIModel(t, "get", "/reports/"+strconv.Itoa(reportID), 200, &getReport)
	assert.EqualValues(t, report.FloorID, getReport.FloorID)
}

func TestListReport(t *testing.T) {
	data := Map{}

	var getReports []Report
	testAPIModelWithQuery(t, "get", "/reports", 200, &getReports, data)
	log.Printf("getReports: %+v\n", getReports)

	data = Map{"range": 1}
	testAPIModelWithQuery(t, "get", "/reports", 200, &getReports, data)
	log.Printf("getReports: %+v\n", getReports)

	data = Map{"range": 2}
	testAPIModelWithQuery(t, "get", "/reports", 200, &getReports, data)
	log.Printf("getReports: %+v\n", getReports)
}

func TestAddReport(t *testing.T) {
	data := Map{"floor_id": REPORT_FLOOR_BASE_ID + 14, "reason": "123456789"}

	testAPI(t, "post", "/reports", 204, data)
}

func TestDeleteReport(t *testing.T) {
	reportID := REPORT_BASE_ID + 7
	var getReport Report
	testAPI(t, "delete", "/reports/"+strconv.Itoa(reportID), 200)

	DB.First(&getReport, reportID)
	assert.EqualValues(t, true, getReport.Dealt)
}
