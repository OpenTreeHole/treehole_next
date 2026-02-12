package tests

import (
	"strconv"
	"testing"

	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func TestGetDivision(t *testing.T) {
	var divisionPinned = []int{0, 2, 3, 1, largeInt}

	var d Division
	DB.First(&d, 1)
	d.Pinned = divisionPinned
	DB.Save(&d)

	var division Division
	testAPIModel(t, "get", "/api/divisions/1", 200, &division)
	// test pinned order
	respPinned := make([]int, 3)
	for i, p := range division.Holes {
		respPinned[i] = p.ID
	}
	assert.Equal(t, []int{2, 3, 1}, respPinned)
}

func TestListDivision(t *testing.T) {
	// return all divisions
	var length int64
	DB.Table("division").Count(&length)
	resp := testAPIArray(t, "get", "/api/divisions", 200)
	assert.Equal(t, length, int64(len(resp)))
}

func TestAddDivision(t *testing.T) {
	data := Map{"name": "TestAddDivision", "description": "TestAddDivisionDescription"}
	testAPI(t, "post", "/api/divisions", 201, data)

	// duplicate post, return 200 and change nothing
	data["description"] = "another"
	resp := testAPI(t, "post", "/api/divisions", 200, data)
	assert.Equal(t, "TestAddDivisionDescription", resp["description"])
}

func TestModifyDivision(t *testing.T) {
	pinned := []int{3, 2, 5, 1, 4}
	data := Map{"name": "modify", "description": "modify", "pinned": pinned}

	var division Division
	testAPIModel(t, "put", "/api/divisions/1", 200, &division, data)

	// test modify
	assert.Equal(t, "modify", division.Name)
	assert.Equal(t, "modify", division.Description)

	// test pinned order
	respPinned := make([]int, 5)
	for i, d := range division.Holes {
		respPinned[i] = d.ID
	}
	assert.Equal(t, pinned, respPinned)
}

func TestDeleteDivision(t *testing.T) {
	id := 3
	toID := 2

	hole := Hole{BaseHole: BaseHole{DivisionID: id}}
	DB.Create(&hole)
	testAPI(t, "delete", "/api/divisions/"+strconv.Itoa(id), 204, Map{"to": toID})
	testAPI(t, "delete", "/api/divisions/"+strconv.Itoa(id), 204, Map{}) // repeat delete

	// deleted
	var d Division
	result := DB.First(&d, id)
	assert.True(t, result.Error != nil)

	// hole moved
	DB.First(&hole, hole.ID)
	assert.Equal(t, toID, hole.DivisionID)

}

func TestDeleteDivisionDefaultValue(t *testing.T) {
	id := 4
	toID := 1

	// if create hole here, say database lock, pending enquiry
	var hole, getHole Hole
	DB.Where("division_id = ?", id).First(&hole)
	testAPI(t, "delete", "/api/divisions/"+strconv.Itoa(id), 204, Map{})

	// hole moved
	DB.Take(&getHole, hole.ID)
	assert.Equal(t, toID, getHole.DivisionID)

}
