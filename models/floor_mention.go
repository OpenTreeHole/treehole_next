package models

import (
	"github.com/gofiber/fiber/v2"
	"regexp"

	"gorm.io/gorm"

	"treehole_next/utils"
)

type FloorMention struct {
	FloorID   int `json:"floor_id" gorm:"primaryKey"`
	MentionID int `json:"mention_id" gorm:"primaryKey"`
}

func (FloorMention) TableName() string {
	return "floor_mention"
}

var reHole = regexp.MustCompile(`[^#]#(\d+)`)
var reFloor = regexp.MustCompile(`##(\d+)`)

func parseMentionIDs(content string) (holeIDs []int, floorIDs []int, err error) {
	// todo: parse replyTo

	// find mentioned holeIDs
	holeIDsText := reHole.FindAllStringSubmatch(" "+content, -1)
	holeIDs, err = utils.RegText2IntArray(holeIDsText)
	if err != nil {
		return nil, nil, err
	}

	// find mentioned floorIDs
	floorIDsText := reFloor.FindAllStringSubmatch(" "+content, -1)
	floorIDs, err = utils.RegText2IntArray(floorIDsText)
	return holeIDs, floorIDs, err
}

func LoadFloorMentions(c *fiber.Ctx, tx *gorm.DB, content string) (Floors, error) {
	holeIDs, floorIDs, err := parseMentionIDs(content)
	if err != nil {
		return nil, err
	}

	if len(holeIDs) == 0 && len(floorIDs) == 0 {
		return Floors{}, nil
	}
	uid, err := FilterUserID(c)
	if err != nil {
		return nil, err
	}
	base := tx.Model(&Floor{})
	if uid != 0 {
		base = base.Where("floor.id IN (?)", VisibleFloorIDs(tx, uid))
	}
	queryGetHoleFloors := base.Session(&gorm.Session{}).Where("hole_id in ? and ranking = 0", holeIDs)
	queryGetFloors := base.Session(&gorm.Session{}).Where("id in ?", floorIDs)
	mentionFloors := Floors{}
	if len(holeIDs) > 0 && len(floorIDs) > 0 {
		err = tx.Raw(`? UNION ?`, queryGetHoleFloors, queryGetFloors).Scan(&mentionFloors).Error
	} else if len(holeIDs) > 0 {
		err = queryGetHoleFloors.Scan(&mentionFloors).Error
	} else if len(floorIDs) > 0 {
		err = queryGetFloors.Scan(&mentionFloors).Error
	}
	return mentionFloors, err
}
