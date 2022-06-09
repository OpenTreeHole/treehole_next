package tests

import (
	"strconv"
	"strings"
	. "treehole_next/models"
)

var largeInt = 1145141919810

func init() {
	for i := 1; i <= 5; i++ {
		division := Division{
			Name:        strconv.Itoa(i),
			Description: strconv.Itoa(i),
		}
		division.ID = i
		DB.Create(&division)
	}
	holes := make([]Hole, 5)
	for i := 0; i < 5; i++ {
		holes[i] = Hole{
			DivisionID: 1,
		}
	}
	tags := make([]Tag, 5)
	tags[0].Holes = []*Hole{&holes[0], &holes[1], &holes[2]}
	tags[0].Temperature = 5
	tags[1].Holes = []*Hole{&holes[3]}
	tags[2].Holes = []*Hole{&holes[4]}
	tags[2].Temperature = 25
	tags[3].Holes = []*Hole{&holes[1], &holes[0], &holes[2]}
	tags[4].Holes = []*Hole{&holes[2], &holes[3], &holes[4]}
	for i := 0; i < 5; i++ {
		tags[i].Name = strings.Repeat("i", i+1)
	}
	DB.Create(&tags)
}
