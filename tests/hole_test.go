package tests

import (
	"strconv"
	"strings"
	"testing"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

const HOLE_BASE = 21

func init() {
	holes := make([]Hole, 10)
	for i := 0; i < 10; i++ {
		holes[i] = Hole{
			DivisionID: 6,
		}
		// holes[i].ID = HOLE_BASE + i
	}
	tag := Tag{
		Name:        "114",
		Temperature: 15,
	}
	holes[1].Tags = []*Tag{&tag}
	holes[2].Tags = []*Tag{&tag}
	holes[3].Tags = []*Tag{
		{
			Name:        "111",
			Temperature: 23,
		},
		{
			Name:        "222",
			Temperature: 45,
		},
	}
	DB.Create(&holes)
	tag = Tag{Name: "115"}
	DB.Create(&tag)
}

func TestGetHoleInDivision(t *testing.T) {
	var holes []Hole
	var ids, respIDs []int

	DB.Raw("SELECT id FROM hole WHERE division_id = 6 AND hidden = 0 ORDER BY updated_at DESC").Scan(&ids)

	testAPIModel(t, "get", "/api/divisions/6/holes", 200, &holes)
	for _, hole := range holes {
		respIDs = append(respIDs, hole.ID)
	}
	assert.Equal(t, ids, respIDs)

	testAPIModel(t, "get", "/api/divisions/"+strconv.Itoa(largeInt)+"/holes", 200, &holes)        // return empty holes
	testAPI(t, "get", "/api/divisions/"+strings.Repeat(strconv.Itoa(largeInt), 15)+"/holes", 500) // huge divisionID
}

func TestListHolesByTag(t *testing.T) {
	var tag Tag
	DB.Where("name = ?", "114").First(&tag)
	var holes []Hole
	DB.Model(&tag).Association("Holes").Find(&holes)

	var getholes []Hole
	testAPIModel(t, "get", "/api/tags/114/holes", 200, &getholes)
	assert.EqualValues(t, len(holes), len(getholes))

	// empty holes
	testAPIModel(t, "get", "/api/tags/115/holes", 200, &getholes)
	assert.EqualValues(t, Holes{}, getholes)
}

func TestCreateHole(t *testing.T) {
	content := "abcdef"
	tagName := []Map{{"name": "a"}, {"name": "ab"}, {"name": "abc"}}
	data := Map{"content": content, "tags": tagName}
	testAPI(t, "post", "/api/divisions/1/holes", 201, data)
	testAPI(t, "post", "/api/divisions/1/holes", 201, data)

	var holes Holes
	var tag Tag
	DB.Where("name = ?", "abc").First(&tag)
	DB.Model(&tag).Association("Holes").Find(&holes)
	assert.EqualValues(t, 2, len(holes))
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
	DB.Model(&tag).Association("Holes").Find(&holes)
}

func TestModifyHole(t *testing.T) {
	var tag Tag
	DB.Where("name = ?", "111").First(&tag)
	var holes Holes
	DB.Model(&tag).Association("Holes").Find(&holes)

	tagName := []Map{{"name": "d"}, {"name": "de"}, {"name": "def"}, {"name": "defg"}}
	division_id := 5
	data := Map{"tags": tagName, "division_id": division_id}
	testAPI(t, "put", "/api/holes/"+strconv.Itoa(holes[0].ID), 200, data)

	DB.Preload("Tags").Where("id = ?", holes[0].ID).Find(&holes[0])

	var getTagName []Map
	for _, v := range holes[0].Tags {
		getTagName = append(getTagName, Map{"name": v.Name})
	}
	assert.EqualValues(t, tagName, getTagName)

	// default schemas
	testAPI(t, "put", "/api/holes/"+strconv.Itoa(holes[0].ID), 200, Map{})
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
