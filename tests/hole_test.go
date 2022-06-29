package tests

import (
	"strconv"
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
	holes[1].Tags = []*Tag{
		&tag,
	}
	holes[2].Tags = []*Tag{
		&tag,
	}
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
}

func TestGetHoleInDivision(t *testing.T) {
	var holes []Hole
	var ids, respIDs []int

	DB.Raw("SELECT id FROM hole WHERE division_id = 1 ORDER BY updated_at DESC").Scan(&ids)

	testAPIModel(t, "get", "/divisions/1/holes", 200, &holes)
	for _, hole := range holes {
		respIDs = append(respIDs, hole.ID)
	}
	assert.Equal(t, ids, respIDs)

	testAPI(t, "get", "/divisions/0/holes", 404)
	testAPI(t, "get", "/divisions/4/holes", 404)
	testAPI(t, "get", "/divisions/1145141919810/holes", 404)
	testAPI(t, "get", "/divisions/1145141919810114514191981011451419198101145141919810/holes", 404)
}

func TestListHolesByTag(t *testing.T) {
	var tag Tag
	DB.Where("name = ?", "114").First(&tag)
	var holes []Hole
	DB.Model(&tag).Association("Holes").Find(&holes)

	var getholes []Hole
	testAPIModel(t, "get", "/tags/114/holes", 200, &getholes)
	assert.EqualValues(t, len(holes), len(getholes))
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

	tagName := []Map{{"name": "d"}, {"name": "de"}, {"name": "def"}}
	division_id := 5
	data := Map{"tags": tagName, "division_id": division_id}
	testAPI(t, "put", "/holes/"+strconv.Itoa(holes[0].ID), 200, data)

	DB.Where("id = ?", holes[0].ID).Find(&holes[0])
	assert.Equal(t, division_id, holes[0].DivisionID)
}

func TestDeleteHole(t *testing.T) {
	var hole Hole
	DB.Where("id = ?", 10).Find(&hole)
	testAPI(t, "delete", "/holes/"+strconv.Itoa(hole.ID), 204)
	DB.Where("id = ?", 10).Find(&hole)
	assert.Equal(t, true, hole.Hidden)
}
