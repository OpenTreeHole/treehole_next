package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"sync"
	"time"
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

var tagCache struct {
	sync.RWMutex
	data      Tags
	nameIndex map[string]*Tag
	idIndex   map[int]*Tag
}

func LoadAllTags(tx *gorm.DB) error {
	tagCache.data = make(Tags, 0, 10000)
	err := tx.Order("temperature DESC").Find(&tagCache.data).Error
	if err != nil {
		return err
	}
	tagCache.nameIndex = make(map[string]*Tag, 10000)
	tagCache.idIndex = make(map[int]*Tag, 10000)

	for _, tag := range tagCache.data {
		tagCache.nameIndex[tag.Name] = tag
		tagCache.idIndex[tag.ID] = tag
	}

	return utils.SetCache("tags", tagCache.data, 10*time.Minute)
}

func LoadTagsByID(tagIDs []int) (tags Tags) {
	tagCache.RLock()
	defer tagCache.RUnlock()
	for _, tagID := range tagIDs {
		if tag, ok := tagCache.idIndex[tagID]; ok {
			tags = append(tags, utils.ValueCopy(tag))
		}
	}
	return tags
}

func LoadTagByName(name string) (*Tag, error) {
	tagCache.RLock()
	defer tagCache.RUnlock()
	if tag, ok := tagCache.nameIndex[name]; ok {
		return utils.ValueCopy(tag), nil
	} else {
		return nil, gorm.ErrRecordNotFound
	}
}

var TagUpdateChan = make(chan int, 1000)
var tagUpdateIDs = make(map[int]bool)

// UpdateTagTemperature is a timed task
func UpdateTagTemperature(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("task UpdateTagTemperature stopped")
			return
		case tagID := <-TagUpdateChan:
			tagUpdateIDs[tagID] = true
		case <-ticker.C:
			updateTagTemperature()
		}
	}
}

// updateTagCacheBytes should be wrapped in tagCache write lock
// tagCache.Lock() should not be called twice
func updateTagCacheBytes() error {
	fmt.Println("update tag cache")
	defer fmt.Println("update finished")
	tagCache.RLock()
	defer tagCache.RUnlock()
	return utils.SetCache("tags", tagCache.data, 10*time.Minute)
}

func updateTagTemperature() {
	tagIDs := utils.Keys(tagUpdateIDs)
	tagUpdateIDs = make(map[int]bool)
	if len(tagUpdateIDs) == 0 {
		return
	}
	var tags Tags
	err := DB.Find(&tags, tagIDs).Error
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	if len(tags) == 0 {
		return
	}

	tagCache.Lock()
	defer tagCache.Unlock()
	for _, tag := range tags {
		originTag, ok := tagCache.idIndex[tag.ID]
		if !ok {
			newTag := utils.ValueCopy(tag)
			tagCache.data = append(tagCache.data, newTag)
			tagCache.idIndex[tag.ID] = newTag
			tagCache.nameIndex[tag.Name] = newTag
		} else {
			*originTag = *tag
		}
	}

	sort.Slice(tagCache.data, func(i, j int) bool {
		return tagCache.data[i].Temperature > tagCache.data[j].Temperature
	})

	err = updateTagCacheBytes()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
}

func (tags Tags) checkTags() Tags {
	newTags := make(Tags, 0)

	// read lock, concurrence reading
	tagCache.RLock()
	defer tagCache.RUnlock()

	for _, tag := range tags {
		cachedTag, ok := tagCache.nameIndex[tag.Name]
		if ok {
			*tag = *cachedTag
		} else {
			newTags = append(newTags, tag)
		}
	}
	return newTags
}

func (tags Tags) FindOrCreateTags(tx *gorm.DB) error {
	newTags := tags.checkTags()
	defer func() {
		for _, tag := range tags {
			TagUpdateChan <- tag.ID
		}
	}()
	if len(newTags) == 0 {
		return nil
	}

	// write lock
	err := createNewTags(tx, newTags)
	if err != nil {
		return err
	}

	return updateTagCacheBytes()
}

func createNewTags(tx *gorm.DB, newTags Tags) error {
	tagCache.Lock()
	defer tagCache.Unlock()

	// check whether newTags have been inserted
	tagsNeedInsert := make(Tags, 0, len(newTags))
	for _, tag := range newTags {
		if _, ok := tagCache.nameIndex[tag.Name]; !ok {
			tagsNeedInsert = append(tagsNeedInsert, tag)
		}
	}

	err := tx.Create(&tagsNeedInsert).Error
	if err != nil {
		return err
	}

	// update cache and index
	for _, newTag := range tagsNeedInsert {
		storeTag := utils.ValueCopy(newTag)
		tagCache.data = append(tagCache.data, storeTag)
		tagCache.nameIndex[storeTag.Name] = storeTag
		tagCache.idIndex[storeTag.ID] = storeTag
	}
	return nil
}
