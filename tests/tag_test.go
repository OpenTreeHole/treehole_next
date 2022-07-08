package tests

import (
	"strconv"
	"testing"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func init() {
	holes := make([]Hole, 5)
	tags := make([]Tag, 6)
	hole_tags := [][]int{
		{0, 1, 2},
		{3},
		{0, 4},
		{1, 0, 2},
		{2, 3, 4},
		{0, 4},
	} // int[tag_id][hole_id]

	for i := range holes {
		holes[i].DivisionID = 6
	}

	for i := range tags {
		tags[i].Name = strconv.Itoa(i + 1)
		for _, v := range hole_tags[i] {
			tags[i].Holes = append(tags[i].Holes, &holes[v])
		}
	}

	tags[0].Temperature = 5
	tags[2].Temperature = 25
	tags[5].Temperature = 34
	DB.Create(&tags)
}

func TestListTag(t *testing.T) {
	var length int64
	DB.Table("tag").Count(&length)
	resp := testAPIArray(t, "get", "/tags", 200)
	assert.Equal(t, length, int64(len(resp)))
}

func TestGetTag(t *testing.T) {
	id := 3

	var tag Tag
	DB.First(&tag, id)

	var newTag Tag
	testAPIModel(t, "get", "/tags/"+strconv.Itoa(id), 200, &newTag)
	assert.Equalf(t, tag.Name, newTag.Name, "get tag")
}

func TestCreateTag(t *testing.T) {
	data := Map{"name": "name"}
	testAPI(t, "post", "/tags", 201, data)

	// duplicate post, return 200 and change nothing
	testAPI(t, "post", "/tags", 200, data)
}

func TestModifyTag(t *testing.T) {
	id := 3
	data := Map{"name": "another", "temperature": 34}

	testAPI(t, "put", "/tags/"+strconv.Itoa(id), 200, data)

	var tag Tag
	DB.Model(&Tag{}).First(&tag, 3)
	assert.Equalf(t, "another", tag.Name, "modify tag name")
	assert.Equalf(t, 34, tag.Temperature, "modify tag tempeture")
}

func TestDeleteTag(t *testing.T) {

	// Move holes to existed tag
	id := 5
	toName := "6"
	data := Map{"to": toName}
	testAPI(t, "delete", "/tags/"+strconv.Itoa(id), 200, data)
	var tag Tag
	DB.Where("name = ?", toName).First(&tag)
	associationHolesLen := DB.Model(&tag).Association("Holes").Count()
	assert.EqualValuesf(t, 4, associationHolesLen, "move holes")
	assert.EqualValuesf(t, 39, tag.Temperature, "tag Temperature add")
	tag = Tag{}

	if result := DB.First(&tag, id); result.Error == nil {
		assert.Error(t, result.Error, "delete tags")
	}

	// Duplicated delete holes
	testAPI(t, "delete", "/tags/"+strconv.Itoa(id), 404, data)

	// Move holes to new tag
	id = 8
	data["to"] = "iii555"
	testAPI(t, "delete", "/tags/"+strconv.Itoa(id), 404, data)
}
