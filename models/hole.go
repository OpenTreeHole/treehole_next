package models

import (
	"fmt"
	"strings"
	"time"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HoleFloor struct {
	FirstFloor *Floor   `json:"first_floor"`
	LastFloor  *Floor   `json:"last_floor"`
	Floors     []*Floor `json:"floors"`
}

type Hole struct {
	BaseModel
	DivisionID int                `json:"division_id"`
	UserID     int                `json:"-"` // permission
	Tags       []*Tag             `json:"tags" gorm:"many2many:hole_tags"`
	Floors     []Floor            `json:"-"`
	HoleFloor  HoleFloor          `json:"floors" gorm:"-:all"` // return floors
	View       int                `json:"view"`
	Reply      int                `json:"reply"`
	Hidden     bool               `json:"hidden"`
	Mapping    []AnonynameMapping `json:"-"`
}
type Holes []Hole

type HoleTag struct {
	HoleID int `json:"hole_id"`
	TagID  int `json:"tag_id"`
}

func loadTags(holes []*Hole) error {
	holeIDs := make([]int, len(holes))
	for i, hole := range holes {
		holeIDs[i] = hole.ID
		hole.Tags = make([]*Tag, 0, config.Config.TagSize)
	}

	var holeTags []*HoleTag
	result := DB.Raw(`
		SELECT * FROM hole_tags
		WHERE hole_id IN (?)`, holeIDs,
	).Scan(&holeTags)
	if result.Error != nil {
		return result.Error
	}

	mapping := make(map[int][]int)
	tagIDs := make([]int, len(holeTags))
	for i, holeTag := range holeTags {
		mapping[holeTag.HoleID] = append(mapping[holeTag.HoleID], holeTag.TagID)
		tagIDs[i] = holeTag.TagID
	}

	var tags []*Tag
	result = DB.Raw(`
		SELECT * FROM tag
		WHERE id IN (?)`, tagIDs,
	).Scan(&tags)
	if result.Error != nil {
		return result.Error
	}
	if len(tags) == 0 {
		return nil
	}

	tagMap := make(map[int]*Tag)
	for _, tag := range tags {
		tagMap[tag.ID] = tag
	}

	for _, hole := range holes {
		for _, tagID := range mapping[hole.ID] {
			hole.Tags = append(hole.Tags, tagMap[tagID])
		}
	}

	return nil
}

func loadFloors(holes []*Hole) error {
	holeIDs := make([]int, len(holes))
	for i, hole := range holes {
		holeIDs[i] = hole.ID
		hole.HoleFloor.Floors = make([]*Floor, 0, config.Config.Size)
	}

	var floors []*Floor
	result := DB.Raw(`
		SELECT *
		FROM (
			SELECT *, rank() over
			(PARTITION BY hole_id ORDER BY id ASC) AS ranking
			FROM floor
		) AS a 
		WHERE hole_id IN (?) AND ranking <= ?`,
		holeIDs, config.Config.Size,
	).Scan(&floors)
	if result.Error != nil {
		return result.Error
	}
	if len(floors) == 0 {
		return nil
	}

	var index, left, right int
	for _, floor := range floors {
		if floor.HoleID != holes[index].ID {
			if index != 0 { // set floors
				holes[index].HoleFloor.Floors = floors[left:right]
				left = right
			}
			for i, hole := range holes { // update index
				if hole.ID == floor.HoleID {
					index = i
					break
				}
			}
		}
		right++
	}
	holes[index].HoleFloor.Floors = floors[left:right]

	for _, hole := range holes {
		if len(hole.HoleFloor.Floors) == 0 {
			return nil
		}

		// first floor
		hole.HoleFloor.FirstFloor = hole.HoleFloor.Floors[0]

		// last floor
		if hole.Reply <= config.Config.Size {
			hole.HoleFloor.LastFloor = hole.HoleFloor.Floors[len(hole.HoleFloor.Floors)-1]
		} else {
			var floor Floor
			DB.Where("hole_id = ?", hole.ID).Last(&floor)
			hole.HoleFloor.LastFloor = &floor
		}
	}

	return nil
}

func (hole *Hole) Preprocess(c *fiber.Ctx) error {
	holes := []*Hole{hole}

	err := loadFloors(holes)
	if err != nil {
		return err
	}

	err = loadTags(holes)
	if err != nil {
		return err
	}

	return nil
}

func getCache(key string) (*Hole, error) {
	// TODO: cache
	return nil, nil
}
func (holes Holes) Preprocess(c *fiber.Ctx) error {
	notInCache := make([]*Hole, 0, len(holes))

	for i := 0; i < len(holes); i++ {
		hole, err := getCache("key")
		if err != nil {
			return err
		}
		if hole == nil {
			notInCache = append(notInCache, &holes[i])
		} else {
			holes[i] = *hole
		}
	}
	err := loadFloors(notInCache)
	if err != nil {
		return err
	}

	err = loadTags(notInCache)
	if err != nil {
		return err
	}

	return nil
}

func MakeQuerySet(c *fiber.Ctx) *gorm.DB {
	var user User
	_ = user.GetUser(c)
	if user.CheckPermission(P_ADMIN) {
		return DB
	} else {
		return DB.Where("hidden = ?", false)
	}
}

func (holes *Holes) MakeQuerySet(offset time.Time, size int, c *fiber.Ctx) (tx *gorm.DB) {
	return MakeQuerySet(c).
		Where("updated_at < ?", offset).
		Order("updated_at desc").Limit(size)
}

// SetTags sets tags for a hole
func (hole *Hole) SetTags(tx *gorm.DB, clear bool) error {
	if clear {
		// update tag temperature
		var sql string
		if config.Config.Debug {
			sql = `
			UPDATE tag
			SET temperature = temperature - 1 
			WHERE id IN (
				SELECT tag_id FROM hole_tags WHERE hole_id = ?
			)`
		} else {
			sql = `
			UPDATE tag INNER JOIN hole_tags 
			ON tag.id = hole_tags.tag_id 
			SET temperature = temperature - 1 
			WHERE hole_tags.hole_id = ?`
		}
		result := tx.Exec(sql, hole.ID)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Exec("DELETE FROM hole_tags WHERE hole_id = ?", hole.ID)
		if result.Error != nil {
			return result.Error
		}
	}

	// create tags
	tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
	}).Create(&hole.Tags)

	tagIDs := make([]int, len(hole.Tags))
	for i, tag := range hole.Tags {
		tagIDs[i] = tag.ID
	}

	// create associations
	var builder strings.Builder
	if config.Config.Debug {
		builder.WriteString("INSERT INTO")
	} else {
		builder.WriteString("INSERT IGNORE INTO")
	}
	builder.WriteString(" hole_tags (hole_id, tag_id) VALUES ")
	for i, tagID := range tagIDs {
		builder.WriteString(fmt.Sprintf("(%d, %d)", hole.ID, tagID))
		if i != len(tagIDs)-1 {
			builder.WriteString(",")
		}
	}
	if config.Config.Debug {
		builder.WriteString(" ON CONFLICT DO NOTHING")
	}
	result := tx.Exec(builder.String())
	if result.Error != nil {
		return result.Error
	}

	// update tag temperature
	result = tx.Exec(`
		UPDATE tag 
		SET temperature = temperature + 1 
		WHERE id IN (?)`,
		tagIDs,
	)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (hole *Hole) Create(c *fiber.Ctx, content *string, specialTag *string, db ...*gorm.DB) error {
	var tx *gorm.DB
	if len(db) > 0 {
		tx = db[0]
	} else {
		tx = DB
	}

	// permission
	var user User
	err := user.GetUser(c)
	if err != nil {
		return err
	}
	if user.BanDivision[hole.DivisionID] ||
		*specialTag != "" && !user.CheckPermission(P_OPERATOR) {
		return utils.Forbidden()
	}
	hole.UserID = user.ID

	err = tx.Transaction(func(tx *gorm.DB) error {
		// Create hole
		result := tx.Omit("Tags").Create(hole) // tags are created in AfterCreate hook
		if result.Error != nil {
			return result.Error
		}

		// Bind and Create floor
		floor := Floor{
			HoleID:     hole.ID,
			Content:    *content,
			UserID:     hole.UserID,
			SpecialTag: *specialTag,
			IsMe:       true,
		}
		err := floor.Create(c, tx)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (hole *Hole) AfterCreate(tx *gorm.DB) (err error) {
	err = hole.SetTags(tx, false)
	if err != nil {
		return err
	}

	return nil
}
