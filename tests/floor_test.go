package tests

import (
	"strconv"
	"strings"
	"testing"
	. "treehole_next/config"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func init() {
	holes := make([]Hole, 5)
	for i := 0; i < 5; i++ {
		holes[i] = Hole{
			DivisionID: 7,
		}
	}
	for i := 1; i <= 50; i++ {
		holes[0].Floors = append(holes[0].Floors, Floor{Content: strings.Repeat("1", i)})
	}
	holes[0].Floors[10].Mention = []Floor{
		{HoleID: 102},
		{HoleID: 304},
	}
	holes[0].Floors[11].Mention = []Floor{
		{HoleID: 506},
		{HoleID: 708},
	}
	holes[1].Floors = append(holes[1].Floors, Floor{Content: "123456789"})
	holes[2].Floors = append(holes[2].Floors, Floor{Content: "123456789"})
	holes[3].Floors = append(holes[3].Floors, Floor{Content: "123456789", Like: 5})
	holes[4].Floors = append(holes[3].Floors, Floor{Content: "123456789"})
	DB.Create(&holes)
}

func TestListFloorsInAHole(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Limit(1).Find(&hole)
	var floors []Floor
	testAPIModel(t, "get", "/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors)
	assert.EqualValues(t, Config.Size, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, "1", floors[0].Content)
	}

	// size
	size := 15
	data := Map{"size": size}
	testAPIModelWithQuery(t, "get", "/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors, data)
	assert.EqualValues(t, size, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, "1", floors[0].Content)
	}

	// big size
	size = 38
	data = Map{"size": size}
	testAPIModelWithQuery(t, "get", "/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors, data)
	assert.EqualValues(t, Config.MaxSize, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, "1", floors[0].Content)
	}

	// offset
	offset := 7
	data = Map{"offset": offset}
	testAPIModelWithQuery(t, "get", "/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors, data)
	assert.EqualValues(t, Config.Size, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, strings.Repeat("1", offset+1), floors[0].Content)
	}
}

func TestListFloorsOld(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).First(&hole)
	data := Map{"hole_id": hole.ID}
	var floors []Floor
	testAPIModelWithQuery(t, "get", "/floors", 200, &floors, data)
	assert.EqualValues(t, 50, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, "1", floors[0].Content)
	}
}

func TestGetFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).First(&hole)
	var floor Floor
	DB.Where("hole_id = ?", hole.ID).First(&floor)
	var getfloor Floor
	testAPIModel(t, "get", "/floors/"+strconv.Itoa(floor.ID), 200, &getfloor)
	assert.EqualValues(t, floor.Content, getfloor.Content)

	testAPIModel(t, "get", "/floors/"+strconv.Itoa(largeInt), 404, &getfloor)
}

func TestCreateFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(1).First(&hole)
	content := "123"
	data := Map{"content": content}
	var getfloor Floor
	testAPIModel(t, "post", "/holes/"+strconv.Itoa(hole.ID)+"/floors", 201, &getfloor, data)
	assert.EqualValues(t, content, getfloor.Content)

	var floors []Floor
	DB.Where("hole_id = ?", hole.ID).Find(&floors)
	assert.EqualValues(t, 2, len(floors))

	testAPIModel(t, "post", "/holes/"+strconv.Itoa(123456)+"/floors", 201, &getfloor, data)
}

func TestCreateFloorOld(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(2).First(&hole)
	content := "1234"
	data := Map{"hole_id": hole.ID, "content": content}
	var getfloor Floor
	testAPIModel(t, "post", "/floors", 201, &getfloor, data)
	assert.EqualValues(t, content, getfloor.Content)

	var floors []Floor
	DB.Where("hole_id = ?", hole.ID).Find(&floors)
	assert.EqualValues(t, 2, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, content, floors[1].Content)
	}

	testAPIModel(t, "post", "/holes/"+strconv.Itoa(123456)+"/floors", 201, &getfloor, data)
}

func TestModifyFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(3).First(&hole)
	var floor Floor
	DB.Where("hole_id = ?", hole.ID).First(&floor)
	content := "12341234"
	data := Map{"content": content}
	var getfloor Floor

	// modify content
	testAPI(t, "put", "/floors/"+strconv.Itoa(floor.ID), 200, data)

	DB.Find(&getfloor, floor.ID)
	assert.EqualValues(t, content, getfloor.Content)

	// modify like delete
	data = Map{"like_int": -1}
	testAPI(t, "put", "/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getfloor, floor.ID)
	assert.EqualValues(t, 4, getfloor.Like)

	// modify like add
	data = Map{"like_int": 1}
	testAPI(t, "put", "/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getfloor, floor.ID)
	assert.EqualValues(t, 5, getfloor.Like)

	// modify like add old
	data = Map{"like": "add"}
	testAPI(t, "put", "/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getfloor, floor.ID)
	assert.EqualValues(t, 6, getfloor.Like)

	// modify like reset
	data = Map{"like_int": 0}
	testAPI(t, "put", "/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getfloor, floor.ID)
	assert.EqualValues(t, 0, getfloor.Like)
}

func TestDeleteFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(4).First(&hole)
	var floor Floor
	DB.Where("hole_id = ?", hole.ID).First(&floor)
	content := "1234567"
	data := Map{"delete_reason": content}

	testAPI(t, "delete", "/floors/"+strconv.Itoa(floor.ID), 200, data)

	DB.First(&floor, floor.ID)
	assert.EqualValues(t, true, floor.Deleted)
	var floorHistory FloorHistory
	DB.Where("floor_id = ?", floor.ID).First(&floorHistory)
	assert.EqualValues(t, content, floorHistory.Reason)
}
