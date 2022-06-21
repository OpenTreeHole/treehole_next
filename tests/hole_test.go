package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
	. "treehole_next/models"
)

func init() {
	holes := make([]Hole, 5)
	for i := 0; i < 5; i++ {
		holes[i] = Hole{
			DivisionID: 2,
		}
	}
	DB.Create(&holes)
}

func TestGetHoleInDivision(t *testing.T) {
	var holes []Hole
	var ids, respIDs []int

	DB.Raw("SELECT id FROM hole WHERE division_id = 1").Scan(&ids)

	testAPIModel(t, "get", "/divisions/1/holes", 200, &holes)

	for _, hole := range holes {
		respIDs = append(respIDs, hole.ID)
	}
	assert.Equal(t, ids, respIDs)
}
