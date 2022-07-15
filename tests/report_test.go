package tests

import (
	"fmt"
	"strconv"
	"testing"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func init() {
	floors := make([]Floor, 20)
	for i := range floors {
		floors[i].ID = 101 + i
	}
	reports := make([]Report, 10)
	for i := range reports {
		reports[i].ID = i + 1
		reports[i].FloorID = 101 + i
		if i < 5 {
			reports[i].Dealt = true
		}
	}

	DB.Create(&floors)
	DB.Create(&reports)
}

func TestGetReport(t *testing.T) {
	reportID := 1
	var report Report
	DB.First(&report, reportID)

	var getReport Report
	testAPIModel(t, "get", "/reports/"+strconv.Itoa(reportID), 200, &getReport)
	assert.EqualValues(t, report.FloorID, getReport.FloorID)
}

func TestListReport(t *testing.T) {
	data := Map{"order_by": "id"}

	var getReports []Report
	testAPIModelWithQuery(t, "get", "/reports", 200, &getReports, data)
	fmt.Printf("getReports: %+v\n", getReports)

	data = Map{"order_by": "id", "range": 1}
	testAPIModelWithQuery(t, "get", "/reports", 200, &getReports, data)
	fmt.Printf("getReports: %+v\n", getReports)

	data = Map{"order_by": "id", "range": 2}
	testAPIModelWithQuery(t, "get", "/reports", 200, &getReports, data)
	fmt.Printf("getReports: %+v\n", getReports)
}

func TestAddReport(t *testing.T) {
	data := Map{"floor_id": 115, "reason": "123456789"}

	var getReport Report
	testAPIModel(t, "post", "/reports", 200, &getReport, data)
	assert.EqualValues(t, 115, getReport.FloorID)
}

func TestDeleteReport(t *testing.T) {
	reportID := 8
	var getReport Report
	testAPI(t, "delete", "/reports/"+strconv.Itoa(reportID), 200)

	DB.First(&getReport, reportID)
	assert.EqualValues(t, true, getReport.Dealt)
}
