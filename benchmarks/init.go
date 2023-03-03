package benchmarks

import (
	"fmt"
	"gorm.io/gorm/logger"
	"math/rand"
	"strings"
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

	divisions := make(Divisions, 0, DIVISION_MAX)
	tags := make(Tags, 0, TAG_MAX)
	holes := make(Holes, 0, HOLE_MAX)
	floors := make(Floors, 0, FLOOR_MAX)

	for i := 0; i < DIVISION_MAX; i++ {
		divisions = append(divisions, &Division{
			ID:          i + 1,
			Name:        strings.Repeat("d", i+1),
			Description: strings.Repeat("dd", i+1),
		})
	}

	for i := 0; i < TAG_MAX; i++ {
		content := fmt.Sprintf("%v", rand.Uint64())
		tags = append(tags, &Tag{
			ID:   i + 1,
			Name: content,
		})
	}

	for i := 0; i < HOLE_MAX; i++ {
		generateTag := func() Tags {
			nowTags := make(Tags, rand.Intn(10))
			for i := range nowTags {
				nowTags[i] = tags[rand.Intn(TAG_MAX)]
			}
			return nowTags
		}
		holes = append(holes, &Hole{
			ID:         i + 1,
			UserID:     1,
			DivisionID: rand.Intn(DIVISION_MAX) + 1,
			Tags:       generateTag(),
		})
	}

	for i := 0; i < FLOOR_MAX; i++ {
		content := fmt.Sprintf("%v", rand.Uint64())
		generateMention := func() Floors {
			floorMentions := make(Floors, 0, rand.Intn(10))
			for j := range floorMentions {
				floorMentions[j] = &Floor{ID: rand.Intn(FLOOR_MAX) + 1}
			}
			return floorMentions
		}
		floors = append(floors, &Floor{
			ID:        i + 1,
			Content:   strings.Repeat(content, rand.Intn(2)),
			Anonyname: utils.GenerateName([]string{}),
			HoleID:    rand.Intn(HOLE_MAX) + 1,
			Mention:   generateMention(),
		})
		holes[floors[i].HoleID-1].Reply += 1
	}

	var err error
	err = DB.Create(divisions).Error
	if err != nil {
		panic(err)
	}
	err = DB.Create(tags).Error
	if err != nil {
		panic(err)
	}
	err = DB.Create(holes).Error
	if err != nil {
		panic(err)
	}
	err = DB.Create(floors).Error
	if err != nil {
		panic(err)
	}
}
