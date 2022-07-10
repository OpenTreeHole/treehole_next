package tests

import (
	"strconv"
	"strings"
	"testing"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func init() {
	holes := make([]Hole, 5)
	for i := 0; i < 5; i++ {
		holes[i] = Hole{
			DivisionID: 2,
		}
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

	DB.Raw("SELECT id FROM hole WHERE division_id = 1 AND hidden = 0 ORDER BY updated_at DESC").Scan(&ids)

	testAPIModel(t, "get", "/divisions/1/holes", 200, &holes)
	for _, hole := range holes {
		respIDs = append(respIDs, hole.ID)
	}
	assert.Equal(t, ids, respIDs)

	testAPIModel(t, "get", "/divisions/"+strconv.Itoa(largeInt)+"/holes", 200, &holes)        // return empty holes
	testAPI(t, "get", "/divisions/"+strings.Repeat(strconv.Itoa(largeInt), 15)+"/holes", 500) // huge divisionID
}

func TestListHolesByTag(t *testing.T) {
	var tag Tag
	DB.Where("name = ?", "114").First(&tag)
	var holes []Hole
	DB.Model(&tag).Association("Holes").Find(&holes)

	var getholes []Hole
	testAPIModel(t, "get", "/tags/114/holes", 200, &getholes)
	assert.EqualValues(t, len(holes), len(getholes))

	// empty holes
	testAPIModel(t, "get", "/tags/115/holes", 200, &getholes)
	assert.EqualValues(t, Holes{}, getholes)
}

func TestCreateHole(t *testing.T) {
	content := "abcdef"
	tagName := []Map{{"name": "a"}, {"name": "ab"}, {"name": "abc"}}
	data := Map{"content": content, "tags": tagName}
	testAPI(t, "post", "/divisions/1/holes", 201, data)
	testAPI(t, "post", "/divisions/1/holes", 201, data)

	var holes Holes
	var tag Tag
	DB.Where("name = ?", "abc").First(&tag)
	DB.Model(&tag).Association("Holes").Find(&holes)
	holes.Preprocess()
	assert.Equal(t, "abcdef", holes[0].HoleFloor.FirstFloor.Content)
}

func TestCreateHoleOld(t *testing.T) {
	content := "abcdef"
	tagName := []Map{{"name": "d"}, {"name": "de"}, {"name": "def"}}
	division_id := 1
	data := Map{"content": content, "tags": tagName, "division_id": division_id}
	testAPI(t, "post", "/holes", 201, data)
	tagName = []Map{{"name": "abc"}, {"name": "defg"}, {"name": "de"}}
	data = Map{"content": content, "tags": tagName, "division_id": division_id}
	testAPI(t, "post", "/holes", 201, data)

	var holes Holes
	var tag Tag
	DB.Where("name = ?", "def").First(&tag)
	DB.Model(&tag).Association("Holes").Find(&holes)
	holes.Preprocess()
	assert.Equal(t, "abcdef", holes[0].HoleFloor.FirstFloor.Content)
}

func TestModifyHole(t *testing.T) {
	var tag Tag
	DB.Where("name = ?", "111").First(&tag)
	var holes Holes
	DB.Model(&tag).Association("Holes").Find(&holes)

	tagName := []Map{{"name": "d"}, {"name": "de"}, {"name": "def"}, {"name": "defg"}}
	division_id := 5
	data := Map{"tags": tagName, "division_id": division_id}
	testAPI(t, "put", "/holes/"+strconv.Itoa(holes[0].ID), 200, data)

	DB.Where("id = ?", holes[0].ID).Find(&holes[0])
	holes[0].Preprocess()
	assert.Equal(t, division_id, holes[0].DivisionID)
	var getTagName []Map
	for _, v := range holes[0].Tags {
		getTagName = append(getTagName, Map{"name": v.Name})
	}
	assert.EqualValues(t, tagName, getTagName)

	// default schemas
	data = Map{}
	testAPI(t, "put", "/holes/"+strconv.Itoa(holes[0].ID), 200, data)
	DB.Where("id = ?", holes[0].ID).Find(&holes[0])
	assert.Equal(t, division_id, holes[0].DivisionID)
}

func TestDeleteHole(t *testing.T) {
	var hole Hole
	holeID := 10
	testAPI(t, "delete", "/holes/"+strconv.Itoa(holeID), 204)
	testAPI(t, "delete", "/holes/"+strconv.Itoa(largeInt), 404)
	DB.Where("id = ?", 10).Find(&hole)
	assert.Equal(t, true, hole.Hidden)
}