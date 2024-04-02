package models

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"treehole_next/config"
	"treehole_next/utils"
)

type Hole struct {
	/// saved fields
	ID        int            `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"time_created" gorm:"not null;index:idx_hole_div_cre,priority:2,sort:desc"`
	UpdatedAt time.Time      `json:"time_updated" gorm:"not null;index:idx_hole_div_upd,priority:2,sort:desc"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	/// base info

	// 浏览量
	View int `json:"view" gorm:"not null;default:0"`

	// 回复量（即该洞下 floor 的数量 - 1）
	Reply int `json:"reply" gorm:"not null;default:0"`

	// 是否隐藏，隐藏的洞用户不可见，管理员可见
	Hidden bool `json:"hidden" gorm:"not null;default:false"`

	// 锁定帖子，如果锁定则非管理员无法发帖，也无法修改已有发帖
	Locked bool `json:"locked" gorm:"not null;default:false"`

	Good bool `json:"good" gorm:"not null;default:false"`

	NoPurge bool `json:"no_purge" gorm:"not null;default:false"`

	/// association info, should add foreign key

	// 所属 division 的 id
	DivisionID int `json:"division_id" gorm:"not null;index:idx_hole_div_upd,priority:1;index:idx_hole_div_cre,priority:1"`

	// 洞主 id，管理员不可见
	UserID int `json:"-" gorm:"not null"`

	// tag 列表；不超过 10 个
	Tags Tags `json:"tags" gorm:"many2many:hole_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// 楼层列表
	Floors Floors `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// 匿名映射表
	Mapping Users `json:"-" gorm:"many2many:anonyname_mapping;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// UserSubscriptionHoles
	UserSubscription Users `json:"-" gorm:"many2many:user_subscription;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	/// generated field

	// 兼容旧版 id
	HoleID int `json:"hole_id" gorm:"-:all"`

	// 返回给前端的楼层列表，包括首楼、尾楼和预加载的前 n 个楼层
	HoleFloor struct {
		FirstFloor *Floor `json:"first_floor"` // 首楼
		LastFloor  *Floor `json:"last_floor"`  // 尾楼
		Floors     Floors `json:"prefetch"`    // 预加载的楼层
	} `json:"floors" gorm:"-:all"`
}

func (hole *Hole) GetID() int {
	return hole.ID
}

func (hole *Hole) CacheName() string {
	return fmt.Sprintf("hole_%d", hole.ID)
}

type Holes []*Hole

func IsHolesExist(tx *gorm.DB, holeID []int) bool {
	var num int64
	tx.Model(&Hole{}).Where("id in ?", holeID).Count(&num)
	return num == int64(len(holeID))
}

/**************
	get hole methods
 *******************/

const HoleCacheExpire = time.Minute * 10

func loadTags(holes Holes) (err error) {
	if len(holes) == 0 {
		return nil
	}
	holeIDs := utils.Models2IDSlice(holes)
	for _, hole := range holes {
		hole.Tags = Tags{}
	}

	var holeTags HoleTags
	err = DB.Where("hole_id in ?", holeIDs).Find(&holeTags).Error
	if err != nil {
		return err
	}

	mapping := make(map[int][]int)
	tagIDs := make(map[int]bool)
	for _, holeTag := range holeTags {
		mapping[holeTag.HoleID] = append(mapping[holeTag.HoleID], holeTag.TagID)
		tagIDs[holeTag.TagID] = true
	}

	var tags Tags
	err = DB.Where("id in ?", utils.Keys(tagIDs)).Find(&tags).Error
	if err != nil {
		return err
	}

	tagMap := make(map[int]*Tag)
	for _, tag := range tags {
		// remove sensitive tags
		if !tag.Sensitive() {
			tagMap[tag.ID] = tag
		}
	}

	for _, hole := range holes {
		for _, tagID := range mapping[hole.ID] {
			hole.Tags = append(hole.Tags, tagMap[tagID])
		}
	}

	return nil
}

func loadFloors(holes Holes) error {
	if len(holes) == 0 {
		return nil
	}
	holeIDs := utils.Models2IDSlice(holes)

	// load all floors with holeIDs and ranking < HoleFloorSize or the last floor
	// sorted by hole_id asc first and ranking asc second
	var floors Floors
	err := DB.
		Raw(
			// using file sort
			`SELECT * FROM (? UNION ?) f ORDER BY hole_id, ranking`,
			// use index(idx_hole_ranking), type range, use MRR
			DB.Model(&Floor{}).Where("hole_id in ? and ranking < ?", holeIDs, config.Config.HoleFloorSize),

			// UNION, remove duplications
			// use index(idx_hole_ranking), type eq_ref
			DB.Model(&Floor{}).Where(
				"(hole_id, ranking) in (?)",
				// use index(PRIMARY), type range
				DB.Model(&Hole{}).Select("id", "reply").Where("id in ?", holeIDs),
			),
		).Scan(&floors).Error
	if err != nil {
		return err
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
		if floor.HoleID != holes[index].ID {
			holes[index].Floors = floors[left:right]
			left = right
			index = slices.IndexFunc(holes, func(hole *Hole) bool {
				return hole.ID == floor.HoleID
			})
		}
		right++
	}
	holes[index].Floors = floors[left:right]

	for _, hole := range holes {
		hole.SetHoleFloor()
	}

	return nil
}

func (hole *Hole) Preprocess(c *fiber.Ctx) error {
	return Holes{hole}.Preprocess(c)
}

func (holes Holes) Preprocess(c *fiber.Ctx) error {
	notInCache := make(Holes, 0, len(holes))

	for _, hole := range holes {
		cachedHole := new(Hole)
		ok := utils.GetCache(hole.CacheName(), &cachedHole)
		if !ok {
			notInCache = append(notInCache, hole)
		} else {
			*hole = *cachedHole
		}
	}

	if len(notInCache) > 0 {
		err := UpdateHoleCache(notInCache)
		if err != nil {
			return err
		}
	}

	// preprocess floors after load from hole cache
	floors := make(Floors, 0)
	for _, hole := range holes {
		hole.SetHoleFloor()
		floors = append(floors, hole.Floors...)

		// hide hole if first floor is sensitive
		if hole.HoleFloor.FirstFloor.Sensitive() {
			hole.Hidden = true
		}
	}
	err := floors.Preprocess(c)
	if err != nil {
		return err
	}

	//user, err := GetUser(c)
	//if err != nil {
	//	return err
	//}
	//
	//// only admin can see hole is hidden
	//if !user.IsAdmin {
	//	for i, hole := range holes {
	//		err = holes[i].Floors.Preprocess(c)
	//		if err != nil {
	//			return err
	//		}
	//		hole.Hidden = false
	//	}
	//}

	return nil
}

func UpdateHoleCache(holes Holes) (err error) {
	err = loadFloors(holes)
	if err != nil {
		return
	}

	err = loadTags(holes)
	if err != nil {
		return
	}

	for _, hole := range holes {
		err = utils.SetCache(hole.CacheName(), hole, HoleCacheExpire)
		if err != nil {
			return
		}
	}
	return
}

func MakeQuerySet(c *fiber.Ctx) (*gorm.DB, error) {
	user, err := GetUser(c)
	if err != nil {
		return nil, err
	}
	if user.IsAdmin {
		return DB, err
	} else {
		return DB.Where("hidden = ?", false), err
	}
}

func (holes Holes) MakeQuerySet(offset common.CustomTime, size int, order string, c *fiber.Ctx) (*gorm.DB, error) {
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

// SetHoleFloor godoc
// set hole.HoleFloor from hole.Floors or hole.HoleFloor.Floors
// if Floors is not empty, set HoleFloor.Floors from Floors, in case loading from database
// if Floors is empty, set HoleFloor.Floors from HoleFloor.Floors, in case loading from cache
func (hole *Hole) SetHoleFloor() {
	if len(hole.Floors) != 0 {
		holeFloorSize := len(hole.Floors)

		hole.HoleFloor.FirstFloor = hole.Floors[0]
		hole.HoleFloor.LastFloor = hole.Floors[holeFloorSize-1]
		if holeFloorSize <= config.Config.HoleFloorSize {
			hole.HoleFloor.Floors = hole.Floors
		} else {
			hole.HoleFloor.Floors = hole.Floors[0 : holeFloorSize-1]
		}
	} else if len(hole.HoleFloor.Floors) != 0 {
		holeFloorSize := len(hole.Floors)

		hole.HoleFloor.FirstFloor = hole.HoleFloor.Floors[0]
		hole.HoleFloor.LastFloor = hole.HoleFloor.Floors[holeFloorSize-1]
		hole.Floors = hole.HoleFloor.Floors
	}

	//for _, floor := range hole.HoleFloor.Floors {
	//	floor.SetDefaults(c)
	//}
	//hole.HoleFloor.FirstFloor.SetDefaults(c)
	//if hole.HoleFloor.FirstFloor.Content == "" {
	//	hole.Hidden = true
	//}
	//hole.HoleFloor.LastFloor.SetDefaults(c)
}

func (hole *Hole) Create(tx *gorm.DB, user *User, tagNames []string, c *fiber.Ctx) (err error) {
	// Create hole.Tags, in different sql session
	hole.Tags, err = FindOrCreateTags(tx, user, tagNames)
	if err != nil {
		return err
	}

	var firstFloor = hole.Floors[0]

	// Find floor.Mentions, in different sql session
	firstFloor.Mention, err = LoadFloorMentions(tx, firstFloor.Content)

	err = tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		// Create hole
		err = tx.Omit(clause.Associations).Create(&hole).Error
		if err != nil {
			return err
		}
		firstFloor.HoleID = hole.ID

		// Create hole_tags association only
		err = tx.Omit("Tags.*", "UpdatedAt").Select("Tags").Save(&hole).Error
		if err != nil {
			return err
		}

		// Update tag temperature
		err = tx.Model(&hole.Tags).Update("temperature", gorm.Expr("temperature + 1")).Error
		if err != nil {
			return err
		}

		// New anonyname
		firstFloor.Anonyname, err = NewAnonyname(tx, hole.ID, hole.UserID)
		if err != nil {
			return err
		}

		// Create floor, set floor_mention association in AfterCreate hook
		return tx.Omit(clause.Associations).Create(&firstFloor).Error
	})
	// transaction commit here
	if err != nil {
		return err
	}

	// set hole.HoleFloor
	hole.SetHoleFloor()

	// half preprocess hole.Floor
	err = firstFloor.SetDefaults(c)
	if err != nil {
		return err
	}

	// index
	if !firstFloor.Sensitive() {
		go FloorIndex(FloorModel{
			ID:        firstFloor.ID,
			UpdatedAt: time.Now(),
			Content:   firstFloor.Content,
		})
	} else {
		firstFloor.Content = ""
	}

	// store into cache
	return utils.SetCache(hole.CacheName(), hole, HoleCacheExpire)
}

func (hole *Hole) AfterCreate(_ *gorm.DB) (err error) {
	hole.HoleID = hole.ID
	return nil
}

func (hole *Hole) AfterFind(_ *gorm.DB) (err error) {
	hole.HoleID = hole.ID
	return nil
}
