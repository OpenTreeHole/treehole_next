package tests

import (
	"strconv"
	. "treehole_next/models"
)

var largeInt = 1145141919810

func init() {
	divisionsSize := 5
	holesSize := 5
	tagsSize := 6

	divisions := make([]Division, divisionsSize)
	holes := make([]Hole, holesSize)
	tags := make([]Tag, tagsSize)
	hole_tags := [][]int{
		{0, 1, 2},
		{3},
		{0, 4},
		{1, 0, 2},
		{2, 3, 4},
		{0, 4},
	} // int[tag_id][hole_id]

	for i := 0; i < divisionsSize; i++ {
		divisions[i].Name = strconv.Itoa(i + 1)
		divisions[i].Description = strconv.Itoa(i + 1)
		divisions[i].ID = i + 1
	}

	for i := 0; i < holesSize; i++ {
		holes[i].DivisionID = 1
		holes[i].ID = i + 1
	}

	for i := 0; i < tagsSize; i++ {
		tags[i].ID = i + 1
		tags[i].Name = strconv.Itoa(i + 1)
		for _, v := range hole_tags[i] {
			tags[i].Holes = append(tags[i].Holes, &holes[v])
		}
	}

	tags[0].Temperature = 5
	tags[2].Temperature = 25
	tags[5].Temperature = 34

	DB.Create(&divisions)
	DB.Create(&tags)
	// when create tags, holes auto create 
}
