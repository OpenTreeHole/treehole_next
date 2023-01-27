package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"treehole_next/utils"
)

type FloorMention struct {
	FloorID   int `json:"floor_id" gorm:"primaryKey"`
	MentionID int `json:"mention_id" gorm:"primaryKey"`
}

func (FloorMention) TableName() string {
	return "floor_mention"
}

func newFloorMentions(floorID int, mentionIDs []int) []FloorMention {
	floorMentions := make([]FloorMention, 0, len(mentionIDs))
	for _, mentionID := range mentionIDs {
		floorMentions = append(floorMentions, FloorMention{
			FloorID:   floorID,
			MentionID: mentionID,
		})
	}
	return floorMentions
}

func deleteFloorMentions(tx *gorm.DB, floorID int) error {
	return tx.Where("floor_id = ?", floorID).Delete(&FloorMention{}).Error
}

func insertFloorMentions(tx *gorm.DB, floorMentions []FloorMention) error {
	if len(floorMentions) == 0 {
		return nil
	}
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(floorMentions).Error
}

var reHole = regexp.MustCompile(`[^#]#(\d+)`)
var reFloor = regexp.MustCompile(`##(\d+)`)

func parseFloorMentions(tx *gorm.DB, content string) ([]int, error) {
	// find mention IDs
	holeIDsText := reHole.FindAllStringSubmatch(" "+content, -1)
	holeIds, err := utils.ReText2IntArray(holeIDsText)
	if err != nil {
		return nil, err
	}

	var mentionIDs = make([]int, 0)
	if len(holeIds) != 0 {
		err := tx.
			Raw("SELECT MIN(id) FROM floor WHERE hole_id IN ? GROUP BY hole_id", holeIds).
			Scan(&mentionIDs).Error
		if err != nil {
			return nil, err
		}
	}

	floorIDsText := reFloor.FindAllStringSubmatch(" "+content, -1)
	mentionIDs2, err := utils.ReText2IntArray(floorIDsText)
	if err != nil {
		return nil, err
	}

	mentionIDs = append(mentionIDs, mentionIDs2...)
	return mentionIDs, nil
}
