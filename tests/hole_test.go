package tests

import (
	"strconv"
	"strings"
	"testing"
	. "treehole_next/config"
	. "treehole_next/models"
	"treehole_next/utils"

	"github.com/stretchr/testify/assert"
)

func TestListHoleInADivision(t *testing.T) {
	var holes Holes
	var ids []int

	DB.Raw("SELECT id FROM hole WHERE division_id = 6 AND hidden = 0 ORDER BY updated_at DESC").Scan(&ids)

	testAPIModel(t, "get", "/api/divisions/6/holes", 200, &holes)
	assert.Equal(t, ids[:Config.HoleFloorSize], utils.Models2IDSlice(holes))

	testAPIModel(t, "get", "/api/divisions/"+strconv.Itoa(largeInt)+"/holes", 200, &holes)        // return empty holes
	testAPI(t, "get", "/api/divisions/"+strings.Repeat(strconv.Itoa(largeInt), 15)+"/holes", 500) // huge divisionID
}

func TestListHolesByTag(t *testing.T) {
	var tag Tag
	DB.Where("name = ?", "114").First(&tag)
	var holes Holes
	err := DB.Model(&tag).Association("Holes").Find(&holes)
	if err != nil {
		t.Fatal(err)
	}

	var getHoles Holes
	testAPIModel(t, "get", "/api/tags/114/holes", 200, &getHoles)
	assert.EqualValues(t, len(holes), len(getHoles))

	// empty holes
	testAPIModel(t, "get", "/api/tags/115/holes", 200, &getHoles)
	assert.EqualValues(t, Holes{}, getHoles)
}

func TestCreateHole(t *testing.T) {
	content := "abcdef"
	data := Map{"content": content, "tags": []Map{{"name": "a"}, {"name": "ab"}, {"name": "abc"}}}
	testAPI(t, "post", "/api/divisions/1/holes", 201, data)
	data["tags"] = []Map{{"name": "abcd"}, {"name": "ab"}, {"name": "abc"}} // update temperature or create tag
	testAPI(t, "post", "/api/divisions/1/holes", 201, data)

	tag := Tag{}
	DB.Where("name = ?", "a").First(&tag)
	assert.EqualValues(t, 1, tag.Temperature)
	tag = Tag{}
	DB.Where("name = ?", "abc").First(&tag)
	assert.EqualValues(t, 2, tag.Temperature)
	assert.EqualValues(t, 2, DB.Model(&tag).Association("Holes").Count())

	data = Map{"content": content}
	testAPI(t, "post", "/api/divisions/1/holes", 400, data) // at least one tag

	content = strings.Repeat("~", 15001)
	data = Map{"content": content, "tags": []Map{{"name": "a"}, {"name": "ab"}, {"name": "abc"}}}
	testAPI(t, "post", "/api/divisions/1/holes", 400, data) // data no more than 15000

	tags := make([]Map, 11)
	for i := range tags {
		tags[i] = Map{"name": strconv.Itoa(i)}
	}
	data = Map{"content": "123456789", "tags": tags} // at most 10 tags
	testAPI(t, "post", "/api/divisions/1/holes", 400, data)
}

func TestCreateHoleOld(t *testing.T) {
	content := "abcdef"
	tagName := []Map{{"name": "d"}, {"name": "de"}, {"name": "def"}}
	division_id := 1
	data := Map{"content": content, "tags": tagName, "division_id": division_id}
	testAPI(t, "post", "/api/holes", 201, data)
	tagName = []Map{{"name": "abc"}, {"name": "defg"}, {"name": "de"}}
	data = Map{"content": content, "tags": tagName, "division_id": division_id}
	testAPI(t, "post", "/api/holes", 201, data)

	var holes Holes
	var tag Tag
	DB.Where("name = ?", "def").First(&tag)
	err := DB.Model(&tag).Association("Holes").Find(&holes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestModifyHole(t *testing.T) {
	var tag Tag
	DB.Where("name = ?", "111").First(&tag)
	var holes Holes
	err := DB.Model(&tag).Association("Holes").Find(&holes)
	if err != nil {
		t.Fatal(err)
	}

	tagName := []Map{{"name": "111"}, {"name": "d"}, {"name": "de"}, {"name": "def"}}
	division_id := 5
	data := Map{"tags": tagName, "division_id": division_id}
	testAPI(t, "put", "/api/holes/"+strconv.Itoa(holes[0].ID), 200, data)

	DB.Preload("Tags").Where("id = ?", holes[0].ID).Find(&holes[0])

	var getTagName []Map
	for _, v := range holes[0].Tags {
		getTagName = append(getTagName, Map{"name": v.Name})
	}
	assert.EqualValues(t, tagName, getTagName)
	assert.EqualValues(t, division_id, holes[0].DivisionID)

	// default schemas
	testAPI(t, "put", "/api/holes/"+strconv.Itoa(holes[0].ID), 400, Map{}) // bad request if modify nothing
	DB.Where("id = ?", holes[0].ID).Find(&holes[0])
	assert.Equal(t, division_id, holes[0].DivisionID)
}

func TestDeleteHole(t *testing.T) {
	var hole Hole
	holeID := 10
	testAPI(t, "delete", "/api/holes/"+strconv.Itoa(holeID), 204)
	testAPI(t, "delete", "/api/holes/"+strconv.Itoa(largeInt), 404)
	DB.Where("id = ?", 10).Find(&hole)
	assert.Equal(t, true, hole.Hidden)
}
