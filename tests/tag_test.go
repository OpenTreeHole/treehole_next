package tests

import (
	"fmt"
	"strconv"
	"testing"
	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

func TestListTag(t *testing.T) {
	var length int64
	DB.Table("tag").Count(&length)
	resp := testAPIArray(t, "get", "/tags", 200)
	assert.Equal(t, length, int64(len(resp)))
}

func TestGetTag(t *testing.T) {
	id := 3

	tag := Tag{}
	tag.ID = id
	DB.First(&tag)
	fmt.Println(tag)

	var newTag Tag
	testAPIModel(t, "get", "/tags/"+strconv.Itoa(id), 200, &newTag)
	assert.Equalf(t, tag.Name, newTag.Name, "get tag")
}

func TestCreateTag(t *testing.T) {
	data := Map{"name": "name"}
	testAPI(t, "post", "/tags", 201, data)

	// duplicate post, return 200 and change nothing
	resp := testAPI(t, "post", "/tags", 200, data)
	fmt.Println(resp)
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
	id := 1
	data := Map{"to": "another"}
	testAPI(t, "delete", "/tags/"+strconv.Itoa(id), 200, data)
	var tag Tag
	DB.Where("name = ?", "another").First(&tag)
	associationHolesLen := DB.Model(&tag).Association("Holes").Count()
	assert.EqualValuesf(t, 4, associationHolesLen, "move holes")
	assert.EqualValuesf(t, 39, tag.Temperature, "tag Temperature add")
	tag = Tag{}
	
	if result := DB.First(&tag, id); result.Error != nil {
		assert.Error(t, result.Error, "delete tags")
	}

	// Duplicated delete holes
	testAPI(t, "delete", "/tags/"+strconv.Itoa(id), 404, data)

	// Move holes to new tag
	id = 4
	data["to"] = "iii555"
	testAPI(t, "delete", "/tags/"+strconv.Itoa(id), 404, data)
}
