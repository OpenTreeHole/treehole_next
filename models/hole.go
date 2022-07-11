package models

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	"strings"
	"time"
	"treehole_next/config"

	"gorm.io/gorm"
)

type HoleFloor struct {
	FirstFloor *Floor   `json:"first_floor"`
	LastFloor  *Floor   `json:"last_floor"`
	Floors     *[]Floor `json:"floors"`
}

type Hole struct {
	BaseModel
	DivisionID int                `json:"division_id"`
	Tags       []*Tag             `json:"tags" gorm:"many2many:hole_tags"`
	Floors     []Floor            `json:"-"`
	HoleFloor  HoleFloor          `json:"floors" gorm:"-:all"` // return floors
	View       int                `json:"view"`
	Reply      int                `json:"reply"`
	Hidden     bool               `json:"hidden"`
	Mapping    []AnonynameMapping `json:"-"`
}
type Holes []Hole

func (hole *Hole) LoadTags() error {
	var tags []*Tag
	err := DB.Model(hole).Association("Tags").Find(&tags)
	if err != nil {
		return err
	}
	if tags == nil {
		hole.Tags = []*Tag{}
	} else {
		hole.Tags = tags
	}
	return nil
}

func (hole *Hole) LoadFloors() error {
	// floors
	var floors []Floor
	result := DB.Where("hole_id = ?", hole.ID).Limit(config.Config.Size).Find(&floors)
	hole.HoleFloor.Floors = &floors
	if result.RowsAffected == 0 {
		return nil
	}

	// first floor
	hole.HoleFloor.FirstFloor = &floors[0]

	// last floor
	if hole.Reply <= config.Config.Size {
		hole.HoleFloor.LastFloor = &floors[result.RowsAffected-1]
	} else {
		var floor Floor
		DB.Where("hole_id = ?", hole.ID).Last(&floor)
		hole.HoleFloor.LastFloor = &floor
	}

	return nil
}

func (hole *Hole) Preprocess() error {
	err := hole.LoadTags()
	if err != nil {
		return err
	}

	err = hole.LoadFloors()
	if err != nil {
		return err
	}

	return nil
}

func (holes Holes) Preprocess() error {
	// TODO: cache
	for i := 0; i < len(holes); i++ {
		if err := holes[i].Preprocess(); err != nil {
			return err
		}
	}
	return nil
}

func MakeQuerySet(c *fiber.Ctx) *gorm.DB {
	var user User
	_ = user.GetUser(c)
	if user.IsAdmin {
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

func (hole *Hole) Create(c *fiber.Ctx, content *string, db ...*gorm.DB) error {
	var tx *gorm.DB
	if len(db) > 0 {
		tx = db[0]
	} else {
		tx = DB
	}

	err := tx.Transaction(func(tx *gorm.DB) error {
		// Create hole
		result := tx.Omit("Tags").Create(hole) // tags are created in AfterCreate hook
		if result.Error != nil {
			return result.Error
		}

		// Bind and Create floor
		floor := Floor{
			HoleID:  hole.ID,
			Content: *content,
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
