package benchmarks

import (
	"fmt"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"math/rand"
	"strings"
	"time"
	. "treehole_next/models"
	"treehole_next/utils"
)

const (
	DIVISION_MAX = 10
	TAG_MAX      = 100
	HOLE_MAX     = 100
	FLOOR_MAX    = 1000
)

func init() {
	DB.Logger = logger.Default.LogMode(logger.Silent)

	rand.Seed(time.Now().UnixMicro())
	divisions := make([]Division, 0, DIVISION_MAX)
	tags := make([]Tag, 0, TAG_MAX)
	holes := make([]Hole, 0, HOLE_MAX)
	floors := make([]Floor, 0, FLOOR_MAX)

	for i := 0; i < DIVISION_MAX; i++ {
		divisions = append(divisions, Division{
			ID:          i + 1,
			Name:        strings.Repeat("d", i+1),
			Description: strings.Repeat("dd", i+1),
		})
	}

	for i := 0; i < TAG_MAX; i++ {
		content := fmt.Sprintf("%v", rand.Uint64())
		tags = append(tags, Tag{
			ID:   i + 1,
			Name: content,
		})
	}

	for i := 0; i < HOLE_MAX; i++ {
		generateTag := func() []*Tag {
			nowtags := make([]*Tag, rand.Intn(10))
			for i := range nowtags {
				nowtags[i] = &tags[rand.Intn(TAG_MAX)]
			}
			return nowtags
		}
		holes = append(holes, Hole{
			ID:         i + 1,
			UserID:     1,
			DivisionID: rand.Intn(DIVISION_MAX) + 1,
			Tags:       generateTag(),
		})
	}

	for i := 0; i < FLOOR_MAX; i++ {
		content := fmt.Sprintf("%v", rand.Uint64())
		generateMention := func() []Floor {
			floors := make([]Floor, rand.Intn(10))
			for i := range floors {
				floors[i].ID = rand.Intn(FLOOR_MAX) + 1
			}
			return floors
		}
		floors = append(floors, Floor{
			ID:        i + 1,
			Content:   strings.Repeat(content, rand.Intn(2)),
			Anonyname: utils.GenerateName([]string{}),
			HoleID:    rand.Intn(HOLE_MAX) + 1,
			Mention:   generateMention(),
		})
		holes[floors[i].HoleID-1].Reply += 1
	}

	createClause := clause.OnConflict{
		UpdateAll: true,
	}
	result := DB.Clauses(createClause).Create(divisions)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	result = DB.Clauses(createClause).Create(tags)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	result = DB.Clauses(createClause).Create(holes)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	result = DB.Clauses(createClause).Create(floors)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}
