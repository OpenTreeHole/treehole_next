package tests

import (
	"strconv"
	"strings"
	"testing"

	"github.com/goccy/go-json"

	. "treehole_next/config"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func TestListFloorsInAHole(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).First(&hole)
	var floors Floors
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

	// size=0 and offset=0 should return all floors
	// The first test hole has 50 floors (see initTestFloors in tests/init.go)
	data = Map{"size": 0, "offset": 0}
	testAPIModelWithQuery(t, "get", "/api/holes/"+strconv.Itoa(hole.ID)+"/floors", 200, &floors, data)
	assert.EqualValues(t, 50, len(floors))
	if len(floors) != 0 {
		assert.EqualValues(t, "1", floors[0].Content)
	}
}

func TestListFloorsOld(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).First(&hole)
	data := Map{"hole_id": hole.ID}
	var floors Floors
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
	var getFloor Floor
	testAPIModel(t, "get", "/api/floors/"+strconv.Itoa(floor.ID), 200, &getFloor)
	assert.EqualValues(t, floor.Content, getFloor.Content)

	testAPIModel(t, "get", "/api/floors/"+strconv.Itoa(largeInt), 404, &getFloor)
}

func TestCreateFloor(t *testing.T) {
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(1).First(&hole)
	content := "123"
	data := Map{"content": content}
	var getFloor Floor
	testAPIModel(t, "post", "/api/holes/"+strconv.Itoa(hole.ID)+"/floors", 201, &getFloor, data)
	assert.EqualValues(t, content, getFloor.Content)

	var floors Floors
	DB.Where("hole_id = ?", hole.ID).Find(&floors)
	assert.EqualValues(t, 2, len(floors))

	testAPIModel(t, "post", "/api/holes/"+strconv.Itoa(largeInt)+"/floors", 404, &getFloor, data)
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
	var getFloor CreateOLdResponse
	rsp := testCommon(t, "post", "/api/floors", 201, data)
	err := json.Unmarshal(rsp, &getFloor)
	assert.Nilf(t, err, "Unmarshal Failed")
	assert.EqualValues(t, content, getFloor.Data.Content)

	var floors Floors
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

	// test6: fold == nil, fold_v2 == "": do nothing; 无效请求
	data = Map{}
	testAPI(t, "put", "/api/floors/"+strconv.Itoa(floor.ID), 400, data)
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
	assert.EqualValues(t, 0, floor.Dislike)

	// dislike
	for i := 0; i < 15; i++ {
		testAPI(t, "post", "/api/floors/"+strconv.Itoa(floor.ID)+"/like/-1", 200)
	}
	DB.First(&floor, floor.ID)
	assert.EqualValues(t, 0, floor.Like)
	assert.EqualValues(t, 1, floor.Dislike)

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
