package tests

import (
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"

	"treehole_next/config"
	. "treehole_next/models"
)

func init() {
	initTestDivision()
	initTestHoles()
	initTestFloors()
	initTestTags()
	initTestFavorites()
	initTestReports()

	config.Config.OpenSensitiveCheck = false
}

func initTestDivision() {
	divisions := make(Divisions, 10)
	for i := range divisions {
		divisions[i] = &Division{
			ID:          i + 1,
			Name:        strconv.Itoa(i),
			Description: strconv.Itoa(i),
		}
	}
	holes := make(Holes, 10)
	for i := range holes {
		holes[i] = &Hole{
			BaseHole: BaseHole{DivisionID: 1},
		}
	}
	holes[9].DivisionID = 4 // for TestDeleteDivisionDefaultValue
	err := DB.Create(&divisions).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	err = DB.Create(&holes).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func initTestHoles() {
	holes := make(Holes, 10)
	for i := range holes {
		holes[i] = &Hole{
			BaseHole: BaseHole{DivisionID: 6},
		}
	}
	tag := Tag{Name: "114", Temperature: 15}
	holes[1].Tags = Tags{&tag}
	holes[2].Tags = Tags{&tag} // here it will insert twice in latest version gorm
	holes[3].Tags = Tags{{Name: "111", Temperature: 23}, {Name: "222", Temperature: 45}}
	err := DB.Create(&holes).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	tag = Tag{Name: "115"}
	err = DB.Create(&tag).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func initTestFloors() {
	holes := make(Holes, 10)
	for i := range holes {
		holes[i] = &Hole{
			BaseHole: BaseHole{DivisionID: 7},
		}
	}
	for i := 1; i <= 50; i++ {
		holes[0].Floors = append(holes[0].Floors, &Floor{Content: strings.Repeat("1", i), Ranking: i - 1})
	}
	holes[0].Floors[10].Mention = Floors{
		{HoleID: 102},
		{HoleID: 304},
	}
	holes[0].Floors[11].Mention = Floors{
		{HoleID: 506},
		{HoleID: 708},
	}
	holes[1].Floors = Floors{{Content: "123456789"}}                                                       // for TestCreate
	holes[2].Floors = Floors{{Content: "123456789"}}                                                       // for TestCreate
	holes[3].Floors = Floors{{Content: "123456789"}}                                                       // for TestModify
	holes[4].Floors = Floors{{Content: "123456789"}}                                                       // for TestModify like
	holes[5].Floors = Floors{{Content: "123456789", UserID: 1}, {Content: "23333", UserID: 5, Ranking: 1}} // for TestDelete
	err := DB.Create(&holes).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func initTestTags() {
	holes := make(Holes, 5)
	tags := make(Tags, 6)
	hole_tags := [][]int{
		{0, 1, 2},
		{3},
		{0, 4},
		{1, 0, 2},
		{2, 3, 4},
		{0, 4},
	} // int[tag_id][hole_id]

	for i := range holes {
		holes[i] = &Hole{
			BaseHole: BaseHole{DivisionID: 8},
		}
	}

	for i := range tags {
		tags[i] = &Tag{Name: strconv.Itoa(i + 1)}
		for _, v := range hole_tags[i] {
			tags[i].Holes = append(tags[i].Holes, holes[v])
		}
	}

	tags[0].Temperature = 5
	tags[2].Temperature = 25
	tags[5].Temperature = 34
	err := DB.Create(&tags).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func initTestFavorites() {
	favoriteGroup := FavoriteGroup{Name: "test", UserID: 1, FavoriteGroupID: 0}
	err := DB.Create(&favoriteGroup).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	userFavorites := make([]UserFavorite, 10)
	for i := range userFavorites {
		userFavorites[i].HoleID = i + 1
		userFavorites[i].UserID = 1
	}
	err = DB.Create(&userFavorites).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

const (
	REPORT_BASE_ID       = 1
	REPORT_FLOOR_BASE_ID = 1001
)

func initTestReports() {
	hole := Hole{BaseHole: BaseHole{ID: 1000}}
	floors := make(Floors, 20)
	for i := range floors {
		floors[i] = &Floor{
			ID:      REPORT_FLOOR_BASE_ID + i,
			HoleID:  1000,
			Ranking: i,
			UserID:  1,
		}
	}
	reports := make([]Report, 10)
	for i := range reports {
		reports[i].ID = REPORT_BASE_ID + i
		reports[i].FloorID = REPORT_FLOOR_BASE_ID + i
		reports[i].UserID = 1
		if i < 5 {
			reports[i].Dealt = true
		}
	}

	err := DB.Create(&hole).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	err = DB.Create(&floors).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	err = DB.Create(&reports).Error
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
