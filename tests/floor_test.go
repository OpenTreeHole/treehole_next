package tests

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	. "treehole_next/config"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func init() {
	holes := make([]Hole, 10)
	for i := 0; i < 10; i++ {
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
	holes[1].Floors = []Floor{{Content: "123456789"}}                                           // for TestCreate
	holes[2].Floors = []Floor{{Content: "123456789"}}                                           // for TestCreate
	holes[3].Floors = []Floor{{Content: "123456789"}}                                           // for TestModify
	holes[4].Floors = []Floor{{Content: "123456789"}}                                           // for TestModify like
	holes[5].Floors = []Floor{{Content: "123456789", UserID: 1}, {Content: "23333", UserID: 5}} // for TestDelete
	DB.Create(&holes)
}

func TestListFloorsInAHole(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).First(&hole)
	var floors []Floor
	testAPIModel(t, "get", "/api/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors)
	assert.EqualValues(t, Config.Size, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, "1", floors[0].Content)
	}

	// size
	size := 38
	data := Map{"size": size}
	testAPIModelWithQuery(t, "get", "/api/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors, data)
	assert.EqualValues(t, size, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, "1", floors[0].Content)
	}

	// offset
	offset := 7
	data = Map{"offset": offset}
	testAPIModelWithQuery(t, "get", "/api/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors, data)
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
	testAPIModelWithQuery(t, "get", "/api/floors", 200, &floors, data)
	assert.EqualValues(t, Config.MaxSize, len(floors))
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
	testAPIModel(t, "get", "/api/floors/"+strconv.Itoa(floor.ID), 200, &getfloor)
	assert.EqualValues(t, floor.Content, getfloor.Content)

	testAPIModel(t, "get", "/api/floors/"+strconv.Itoa(largeInt), 404, &getfloor)
}

func TestCreateFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(1).First(&hole)
	content := "123"
	data := Map{"content": content}
	var getfloor Floor
	testAPIModel(t, "post", "/api/holes/"+strconv.Itoa(hole.ID)+"/floors", 201, &getfloor, data)
	assert.EqualValues(t, content, getfloor.Content)

	var floors []Floor
	DB.Where("hole_id = ?", hole.ID).Find(&floors)
	assert.EqualValues(t, 2, len(floors))

	testAPIModel(t, "post", "/api/holes/"+strconv.Itoa(largeInt)+"/floors", 404, &getfloor, data)
}

func TestCreateFloorOld(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(2).First(&hole)
	content := "1234"
	data := Map{"hole_id": hole.ID, "content": content}
	type CreateOLdResponse struct {
		Data    Floor
		Message string
	}
	var getfloor CreateOLdResponse
	rsp := testCommon(t, "post", "/api/floors", 201, data)
	err := json.Unmarshal(rsp, &getfloor)
	assert.Nilf(t, err, "Unmarshal Failed")
	assert.EqualValues(t, content, getfloor.Data.Content)

	var floors []Floor
	DB.Where("hole_id = ?", hole.ID).Find(&floors)
	assert.EqualValues(t, 2, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, content, floors[1].Content)
	}

	testCommon(t, "post", "/api/holes/"+strconv.Itoa(123456)+"/floors", 404, data)
}

func TestModifyFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(3).First(&hole)
	var floor Floor
	DB.Where("hole_id = ?", hole.ID).First(&floor)
	content := "12341234"
	data := Map{"content": content}
	var getFloor Floor

	// modify content
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)

	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, content, getFloor.Content)

	// modify fold
	// test 1: fold == ["test"], fold_v2 == ""
	data = Map{"fold": []string{"test"}}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, data["fold"].([]string)[0], getFloor.Fold)

	// test2: fold == [], fold_v2 == "": expect reset fold
	data = Map{"fold": []string{}}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, "", getFloor.Fold)

	// test3: fold == [], fold_v2 == "test_test"
	data = Map{"fold_v2": "test_test"}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, data["fold_v2"], getFloor.Fold)

	// test4: fold == [], fold_v2 == "": expect reset fold
	data = Map{"fold": []string{}}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, "", getFloor.Fold)

	// test5: fold == ["test"], fold_v2 == "test_test": expect "test_test", fold_v2 has the priority
	data = Map{"fold": []string{"test"}, "fold_v2": "test_test"}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, "test_test", getFloor.Fold)

	// test6: fold == nil, fold_v2 == "": do nothing
	data = Map{}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, "test_test", getFloor.Fold)

	// modify like add old
	data = Map{"like": "add"}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, 1, getFloor.Like)

	// modify like reset old
	data = Map{"like": "cancel"}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
	DB.Find(&getFloor, floor.ID)
	assert.EqualValues(t, 0, getFloor.Like)
}

func TestModifyFloorLike(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(4).First(&hole)
	var floor Floor
	DB.Where("hole_id = ?", hole.ID).First(&floor)

	// like
	for i := 0; i < 10; i++ {
		testAPI(t, "post", "/api/floors/"+strconv.Itoa(floor.ID)+"/like/1", 200)
	}
	DB.First(&floor, floor.ID)
	assert.EqualValues(t, 1, floor.Like)

	// dislike
	for i := 0; i < 15; i++ {
		testAPI(t, "post", "/api/floors/"+strconv.Itoa(floor.ID)+"/like/-1", 200)
	}
	DB.First(&floor, floor.ID)
	assert.EqualValues(t, -1, floor.Like)

	// reset
	testAPI(t, "post", "/api/floors/"+strconv.Itoa(floor.ID)+"/like/0", 200)
	DB.First(&floor, floor.ID)
	assert.EqualValues(t, 0, floor.Like)
}

func TestDeleteFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(5).First(&hole)
	var floor Floor
	DB.Where("hole_id = ?", hole.ID).First(&floor)
	content := "1234567"
	data := Map{"delete_reason": content}

	testAPI(t, "delete", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)

	DB.First(&floor, floor.ID)
	assert.EqualValues(t, true, floor.Deleted)
	var floorHistory FloorHistory
	DB.Where("floor_id = ?", floor.ID).First(&floorHistory)
	assert.EqualValues(t, content, floorHistory.Reason)

	// permission
	floor = Floor{}
	DB.Where("hole_id = ?", hole.ID).Offset(1).First(&floor)
	testAPI(t, "delete", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)
}
