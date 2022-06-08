package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	. "treehole_next/models"
	"treehole_next/schemas"
)

func init() {
	for i := 1; i <= 5; i++ {
		division := Division{
			Name:        strconv.Itoa(i),
			Description: strconv.Itoa(i),
		}
		division.ID = i
		DB.Create(&division)
	}
	holes := make([]Hole, 5)
	for i := 0; i < 5; i++ {
		holes[i] = Hole{
			DivisionID: 1,
		}
	}
	DB.Create(&holes)
}

func TestGetDivision(t *testing.T) {
	var divisionPinned = []int{0, 2, 3, 1, largeInt}

	var d Division
	DB.First(&d, 1)
	d.Pinned = divisionPinned
	DB.Save(&d)

	var division schemas.DivisionResponse
	testAPIModel(t, "get", "/divisions/1", 200, &division)
	// test pinned order
	respPinned := make([]int, 3)
	for i, p := range division.Pinned {
		respPinned[i] = p.ID
	}
	assert.Equal(t, []int{2, 3, 1}, respPinned)
}

func TestListDivision(t *testing.T) {
	// return all divisions
	var length int64
	DB.Table("division").Count(&length)
	resp := testAPIArray(t, "get", "/divisions", 200)
	assert.Equal(t, length, int64(len(resp)))
}

func TestAddDivision(t *testing.T) {
	data := Map{"name": "name", "description": "description"}
	testAPI(t, "post", "/divisions", 201, data)

	// duplicate post, return 200 and change nothing
	data["description"] = "another"
	resp := testAPI(t, "post", "/divisions", 200, data)
	fmt.Println(resp)
	assert.Equal(t, "description", resp["description"])
}

func TestModifyDivision(t *testing.T) {
	pinned := []int{3, 2, 5, 1, 4}
	data := Map{"name": "modify", "description": "modify", "pinned": pinned}

	var division schemas.DivisionResponse
	testAPIModel(t, "put", "/divisions/1", 200, &division, data)

	// test modify
	assert.Equal(t, "modify", division.Name)
	assert.Equal(t, "modify", division.Description)

	// test pinned order
	respPinned := make([]int, 5)
	for i, d := range division.Pinned {
		respPinned[i] = d.ID
	}
	assert.Equal(t, pinned, respPinned)
}

func TestDeleteDivision(t *testing.T) {
	id := 3
	toID := 2

	hole := Hole{DivisionID: id}
	DB.Create(&hole)
	testAPI(t, "delete", "/divisions/"+strconv.Itoa(id), 204, Map{"to": toID})
	testAPI(t, "delete", "/divisions/"+strconv.Itoa(id), 204) // repeat delete

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

	hole := Hole{DivisionID: id}
	DB.Create(&hole)
	testAPI(t, "delete", "/divisions/"+strconv.Itoa(id), 204)

	// hole moved
	DB.First(&hole, hole.ID)
	assert.Equal(t, toID, hole.DivisionID)

}
