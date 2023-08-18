package models

import (
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

func LoadFloorMentions(tx *gorm.DB, content string) (Floors, error) {
	holeIDs, floorIDs, err := parseMentionIDs(content)
	if err != nil {
		return nil, err
	}

	queryGetHoleFloors := tx.Model(&Floor{}).Where("hole_id in ? and ranking = 0", holeIDs)
	queryGetFloors := tx.Model(&Floor{}).Where("id in ?", floorIDs)
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
