package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"treehole_next/utils"
)

type AnonynameMapping struct {
	HoleID    int    `json:"hole_id" gorm:"primaryKey"`
	UserID    int    `json:"user_id" gorm:"primaryKey"`
	Anonyname string `json:"anonyname" gorm:"size:32"`
}

func NewAnonyname(tx *gorm.DB, holeID, userID int) (string, error) {
	name := utils.NewRandName()
	return name, tx.Create(&AnonynameMapping{
		HoleID:    holeID,
		UserID:    userID,
		Anonyname: name,
	}).Error
}

func FindOrGenerateAnonyname(tx *gorm.DB, holeID, userID int) (string, error) {
	var anonyname string
	err := tx.
		Model(&AnonynameMapping{}).
		Select("anonyname").
		Where("hole_id = ?", holeID).
		Where("user_id = ?", userID).
		Take(&anonyname).Error

	if err == gorm.ErrRecordNotFound {
		var names []string
		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Model(&AnonynameMapping{}).
			Select("anonyname").
			Where("hole_id = ?", holeID).
			Order("anonyname").
			Scan(&names).Error
		if err != nil {
			return "", err
		}

		anonyname = utils.GenerateName(names)
		err = tx.Create(&AnonynameMapping{
			HoleID:    holeID,
			UserID:    userID,
			Anonyname: anonyname,
		}).Error
		if err != nil {
			return anonyname, err
		}
	} else if err != nil {
		return "", err
	}

	return anonyname, nil
}
