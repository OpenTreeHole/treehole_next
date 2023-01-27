package models

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Hole struct {
	ID         int                `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time          `json:"time_created"`
	UpdatedAt  time.Time          `json:"time_updated"`
	HoleID     int                `json:"hole_id" gorm:"-:all"`                                                          // 兼容旧版 id
	DivisionID int                `json:"division_id"`                                                                   // 所属 division 的 id
	UserID     int                `json:"-"`                                                                             // 洞主 id
	Tags       []*Tag             `json:"tags" gorm:"many2many:hole_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // tag 列表
	Floors     []Floor            `json:"-"`                                                                             // 楼层列表
	HoleFloor  HoleFloor          `json:"floors" gorm:"-:all"`                                                           // 返回给前端的楼层列表，包括首楼、尾楼和预加载的前 n 个楼层
	View       int                `json:"view"`                                                                          // 浏览量
	Reply      int                `json:"reply"`                                                                         // 回复量（即该洞下 floor 的数量）
	Hidden     bool               `json:"hidden" gorm:"index"`                                                           // 是否隐藏，隐藏的洞用户不可见，管理员可见
	Mapping    []AnonynameMapping `json:"-"`                                                                             // 匿名映射表
}

func (hole Hole) GetID() int {
	return hole.ID
}

type Holes []Hole

type HoleFloor struct {
	FirstFloor *Floor   `json:"first_floor"` // 首楼
	LastFloor  *Floor   `json:"last_floor"`  // 尾楼
	Floors     []*Floor `json:"prefetch"`    // 预加载的楼层
}

/**************
	get hole methods
 *******************/

const HoleCacheExpire = time.Minute * 10

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
		tag.TagID = tag.ID
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
		hole.HoleFloor.Floors = make([]*Floor, 0, config.Config.HoleFloorSize)
	}

	var floors []*Floor
	result := DB.Raw(`
		SELECT *
		FROM (
			SELECT *, rank() over
			(PARTITION BY hole_id ORDER BY id) AS ranking
			FROM floor
		) AS a 
		WHERE hole_id IN (?) AND ranking <= ?`,
		holeIDs, config.Config.HoleFloorSize,
	).Scan(&floors)
	if result.Error != nil {
		return result.Error
	}
	if len(floors) == 0 {
		return nil
	}

	/*
			Bind floors to hole.
			Note that floor is grouped by hole_id in hole_id asc order
		and hole is in random order, so we have to find hole_id those floors
		belong to both at the beginning and after floor group has changed.
			To bind, we use two pointers. Binding occurs when the floor's hole_id
		has changed, or when the floor is the last floor.
			The complexity is O(m*n), where m is the number of holes and
		n is the number of floors. Given that m is relatively small,
		the complexity is acceptable.
	*/
	var left, right int
	index := slices.IndexFunc(holes, func(hole *Hole) bool {
		return hole.ID == floors[0].HoleID
	})
	for _, floor := range floors {
		floor.SetDefaults()
		if floor.HoleID != holes[index].ID {
			holes[index].HoleFloor.Floors = floors[left:right]
			left = right
			index = slices.IndexFunc(holes, func(hole *Hole) bool {
				return hole.ID == floor.HoleID
			})
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
		// this means all the floors are loaded into hole.HoleFloor.Floors,
		// so we can just get last floor from hole.HoleFloor.Floors
		if hole.Reply < config.Config.HoleFloorSize {
			hole.HoleFloor.LastFloor = hole.HoleFloor.Floors[len(hole.HoleFloor.Floors)-1]
		} else {
			var floor Floor
			DB.Where("hole_id = ?", hole.ID).Last(&floor)
			floor.SetDefaults()
			hole.HoleFloor.LastFloor = &floor
		}
	}

	return nil
}

func (hole *Hole) Preprocess(c *fiber.Ctx) error {
	holes := Holes{*hole}

	err := holes.Preprocess(c)
	if err != nil {
		return err
	}

	*hole = holes[0]

	return nil
}

func (holes Holes) Preprocess(c *fiber.Ctx) error {
	notInCache := make([]*Hole, 0, len(holes))

	for i := 0; i < len(holes); i++ {
		var hole Hole
		ok := utils.GetCache(fmt.Sprintf("hole_%d", holes[i].ID), &hole)
		if !ok {
			notInCache = append(notInCache, &holes[i])
		} else {
			holes[i] = hole
		}
	}

	if len(notInCache) > 0 {
		err := UpdateHoleCache(notInCache)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateHoleCache(notInCache []*Hole) error {
	err := loadFloors(notInCache)
	if err != nil {
		return err
	}

	err = loadTags(notInCache)
	if err != nil {
		return err
	}

	for i := range notInCache {
		notInCache[i].HoleID = notInCache[i].ID
	}

	for i := 0; i < len(notInCache); i++ {
		err = utils.SetCache(
			fmt.Sprintf("hole_%d", notInCache[i].ID),
			notInCache[i],
			HoleCacheExpire,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func MakeQuerySet(c *fiber.Ctx) (*gorm.DB, error) {
	user, err := GetUser(c)
	if err != nil {
		return nil, err
	}
	if perm.CheckPermission(user, perm.Admin) {
		return DB, err
	} else {
		return DB.Where("hidden = ?", false), err
	}
}

func (holes Holes) MakeQuerySet(offset CustomTime, size int, order string, c *fiber.Ctx) (*gorm.DB, error) {
	querySet, err := MakeQuerySet(c)
	if err != nil {
		return nil, err
	}
	if order == "time_created" || order == "created_at" {
		return querySet.
			Where("created_at < ?", offset.Time).
			Order("created_at desc").Limit(size), nil
	} else {
		return querySet.
			Where("updated_at < ?", offset.Time).
			Order("updated_at desc").Limit(size), nil
	}
}

/************************
	create and modify hole methods
 ************************/

// SetTags sets tags for a hole
func (hole *Hole) SetTags(tx *gorm.DB, clear bool) error {
	var err error
	if clear {
		err = tx.Exec(`
			UPDATE tag SET temperature = temperature - 1 
			WHERE id IN ( SELECT tag_id FROM hole_tags WHERE hole_id = ?)`, hole.ID).Error
		if err != nil {
			return err
		}

		err = tx.Exec("DELETE FROM hole_tags WHERE hole_id = ?", hole.ID).Error
		if err != nil {
			return err
		}
	}

	if len(hole.Tags) == 0 {
		return nil
	}

	// create tags
	result := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&hole.Tags)
	if result.Error != nil {
		return result.Error
	}

	// find tags
	tagNames := make([]string, len(hole.Tags))
	for i, tag := range hole.Tags {
		tagNames[i] = tag.Name
	}
	result = tx.Where("name IN (?)", tagNames).Find(&hole.Tags)
	if result.Error != nil {
		return result.Error
	}

	// create associations
	tagIDs := make([]int, len(hole.Tags))
	for i, tag := range hole.Tags {
		tagIDs[i] = tag.ID
	}
	var builder strings.Builder

	if DBType == DBTypeSqlite {
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

	if DBType == DBTypeSqlite {
		builder.WriteString(" ON CONFLICT DO NOTHING")
	}
	result = tx.Exec(builder.String())
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

func (hole *Hole) Create(c *fiber.Ctx, content string, specialTag string, db ...*gorm.DB) error {
	var tx *gorm.DB
	if len(db) > 0 {
		tx = db[0]
	} else {
		tx = DB
	}

	hole.UserID, _ = GetUserID(c)

	return tx.Transaction(func(tx *gorm.DB) error {
		// Create hole
		hole.Reply = -1
		result := tx.Omit("Tags").Create(hole) // tags are created in AfterCreate hook
		if result.Error != nil {
			return result.Error
		}
		hole.Reply = 0

		// Bind and Create floor
		floor := Floor{
			HoleID:     hole.ID,
			Content:    content,
			UserID:     hole.UserID,
			SpecialTag: specialTag,
		}

		// create floor
		err := floor.Create(c, tx)
		if err != nil {
			return err
		}

		// create Favorite
		return UserCreateFavourite(tx, c, false, hole.UserID, []int{hole.ID})
	})
}

func (hole *Hole) AfterCreate(tx *gorm.DB) (err error) {
	hole.HoleID = hole.ID
	return hole.SetTags(tx, false)
}
