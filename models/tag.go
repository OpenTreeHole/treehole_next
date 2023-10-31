package models

import (
	"strings"
	"time"

	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"treehole_next/utils"
)

type Tag struct {
	/// saved fields
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"-" gorm:"not null"`
	UpdatedAt time.Time `json:"-" gorm:"not null"`

	/// base info
	Name        string `json:"name" gorm:"not null;unique;size:32"`
	Temperature int    `json:"temperature" gorm:"not null;default:0"`

	/// association info, should add foreign key
	Holes Holes `json:"-" gorm:"many2many:hole_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	/// generated field
	TagID int `json:"tag_id" gorm:"-:all"`
}

type Tags []*Tag

func (tag *Tag) GetID() int {
	return tag.ID
}

func (tag *Tag) AfterFind(_ *gorm.DB) (err error) {
	tag.TagID = tag.ID
	return nil
}

func (tag *Tag) AfterCreate(_ *gorm.DB) (err error) {
	tag.TagID = tag.ID
	return nil
}

func FindOrCreateTags(tx *gorm.DB, user *User, names []string) (Tags, error) {
	tags := make(Tags, 0)
	err := tx.Where("name in ?", names).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	existTagName := make([]string, 0)
	for _, tag := range tags {
		existTagName = append(existTagName, tag.Name)
	}

	newTags := make(Tags, 0)
	for _, name := range names {
		name = strings.TrimSpace(name)
		if !slices.ContainsFunc(existTagName, func(s string) bool {
			return strings.EqualFold(s, name)
		}) {
			newTags = append(newTags, &Tag{Name: name})
		}
	}

	if len(newTags) == 0 {
		return tags, nil
	}

	for _, tag := range newTags {
		if strings.HasPrefix(tag.Name, "#") {
			if !user.IsAdmin {
				return nil, common.BadRequest("只有管理员才能创建 # 开头的 tag")
			}
		}
		if strings.HasPrefix(tag.Name, "@") {
			if !user.IsAdmin {
				return nil, common.BadRequest("只有管理员才能创建 @ 开头的 tag")
			}
		}
	}

	err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&newTags).Error

	go UpdateTagCache(nil)

	return append(tags, newTags...), err
}

func UpdateTagCache(tags Tags) {
	var err error
	if len(tags) == 0 {
		err := DB.Order("temperature desc").Find(&tags).Error
		if err != nil {
			log.Printf("update tag cache error: %s", err)
		}
	}
	err = utils.SetCache("tags", tags, 10*time.Minute)
	if err != nil {
		log.Printf("update tag cache error: %s", err)
	}
}
