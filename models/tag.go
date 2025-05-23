package models

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"treehole_next/config"
	"treehole_next/utils/sensitive"

	"github.com/gofiber/fiber/v2"

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

	IsZZMG bool `json:"-" gorm:"not null;default:false"`

	/// association info, should add foreign key
	Holes Holes `json:"-" gorm:"many2many:hole_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// auto sensitive check
	IsSensitive bool `json:"-" gorm:"index:idx_tag_actual_sensitive,priority:1"`

	// manual sensitive check
	IsActualSensitive *bool `json:"-" gorm:"index:idx_tag_actual_sensitive,priority:2"`
	/// generated field
	TagID int `json:"tag_id" gorm:"-:all"`

	Nsfw bool `json:"nsfw" gorm:"not null;default:false;index"`
}

type Tags []*Tag

func (tag *Tag) GetID() int {
	return tag.ID
}

func (tag *Tag) AfterFind(_ *gorm.DB) (err error) {
	tag.TagID = tag.ID
	return nil
}

func (tag *Tag) BeforeCreate(_ *gorm.DB) (err error) {
	if len(tag.Name) > 0 && tag.Name[0] == '*' {
		tag.Nsfw = true
	}
	return nil
}

func (tag *Tag) AfterCreate(_ *gorm.DB) (err error) {
	tag.TagID = tag.ID
	return nil
}

func FindOrCreateTags(tx *gorm.DB, user *User, names []string) (Tags, error) {
	tags := make(Tags, 0)
	for i, name := range names {
		names[i] = strings.TrimSpace(name)
	}
	err := tx.Where("name in ?", names).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	existTagNames := make([]string, 0)
	for _, tag := range tags {
		existTagNames = append(existTagNames, tag.Name)
		if !user.IsAdmin {
			if slices.ContainsFunc(config.Config.AdminOnlyTagIds, func(i int) bool {
				return i == tag.ID
			}) {
				return nil, common.Forbidden(fmt.Sprintf("标签 %s 为管理员专用标签", tag.Name))
			}
		}
	}

	newTags := make(Tags, 0)
	for _, name := range names {
		name = strings.TrimSpace(name)
		if !slices.ContainsFunc(existTagNames, func(s string) bool {
			return strings.EqualFold(s, name)
		}) {
			newTags = append(newTags, &Tag{Name: name})
		}
	}

	if len(newTags) == 0 {
		return tags, nil
	}
	for _, tag := range newTags {
		if !user.IsAdmin {
			if len(tag.Name) > 15 && len([]rune(tag.Name)) > 10 {
				return nil, common.BadRequest("标签长度不能超过 10 个字符")
			}
			if strings.HasPrefix(tag.Name, "#") {
				return nil, common.BadRequest("只有管理员才能创建 # 开头的 tag")
			}
			if strings.HasPrefix(tag.Name, "@") {
				return nil, common.BadRequest("只有管理员才能创建 @ 开头的 tag")
			}
			if strings.HasPrefix(tag.Name, "*") {
				return nil, common.BadRequest("只有管理员才能创建 * 开头的 tag")
			}
		}
	}

	var wg sync.WaitGroup
	for _, tag := range newTags {
		wg.Add(1)
		go func(tag *Tag) {
			sensitiveResp, err := sensitive.CheckSensitive(sensitive.ParamsForCheck{
				Content:  tag.Name,
				Id:       time.Now().UnixNano(),
				TypeName: sensitive.TypeTag,
			})
			if err != nil {
				return
			}
			tag.IsSensitive = !sensitiveResp.Pass
			wg.Done()
		}(tag)
	}
	wg.Wait()

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

func (tag *Tag) Preprocess(c *fiber.Ctx) error {
	return Tags{tag}.Preprocess(c)
}

func (tags Tags) Preprocess(_ *fiber.Ctx) error {
	tagIDs := make([]int, len(tags))
	IdTagMapping := make(map[int]*Tag)
	for i, tag := range tags {
		if tags[i].Sensitive() {
			tags[i].Name = ""
		}
		tagIDs[i] = tag.ID
		IdTagMapping[tag.ID] = tags[i]
	}
	return nil
}

func (tag *Tag) Sensitive() bool {
	if tag == nil {
		return false
	}
	if tag.IsActualSensitive != nil {
		return *tag.IsActualSensitive
	}
	return tag.IsSensitive
}
