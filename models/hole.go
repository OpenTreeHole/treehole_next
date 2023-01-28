package models

import (
	"fmt"
	"golang.org/x/exp/slices"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Hole struct {
	/// saved fields
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created" gorm:"not null;index:idx_hole_div_cre,priority:2"`
	UpdatedAt time.Time `json:"time_updated" gorm:"not null;index:idx_hole_div_upd,priority:2"`

	/// base info

	// 浏览量
	View int `json:"view" gorm:"not null;default:0"`

	// 回复量（即该洞下 floor 的数量 - 1）
	Reply int `json:"reply" gorm:"not null;default:0"`

	// 是否隐藏，隐藏的洞用户不可见，管理员可见
	Hidden bool `json:"hidden" gorm:"not null;default:false"`

	/// association info, should add foreign key

	// 所属 division 的 id
	DivisionID int `json:"division_id" gorm:"not null;index:idx_hole_div_upd,priority:1;index:idx_hole_div_cre,priority:1"`

	// 洞主 id，管理员可见
	UserID int `json:"user_id;omitempty" gorm:"not null"`

	// tag 列表；不超过 10 个
	Tags Tags `json:"tags" gorm:"many2many:hole_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// 楼层列表
	Floors Floors `json:"-"`

	// 匿名映射表
	Mapping Users `json:"-" gorm:"many2many:anonyname_mapping"`

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

func (hole *Hole) IDString() string {
	return fmt.Sprintf(fmt.Sprintf("hole_%d", hole.ID))
}

type Holes []*Hole

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
	return Holes{hole}.Preprocess(c)
}

func (holes Holes) Preprocess(_ *fiber.Ctx) error {
	notInCache := make(Holes, 0, len(holes))

	for i := 0; i < len(holes); i++ {
		hole := new(Hole)
		ok := utils.GetCache(hole.IDString(), &hole)
		if !ok {
			notInCache = append(notInCache, holes[i])
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

func UpdateHoleCache(holes Holes) error {
	err := loadFloors(holes)
	if err != nil {
		return err
	}

	err = loadTags(holes)
	if err != nil {
		return err
	}

	for i := range holes {
		holes[i].HoleID = holes[i].ID
	}

	for i := range holes {
		err = utils.SetCache(
			holes[i].IDString(),
			holes[i],
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
	holeTags := make([]HoleTag, 0, len(hole.Tags))
	for _, tag := range hole.Tags {
		holeTags = append(holeTags, HoleTag{hole.ID, tag.ID})
	}
	err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(holeTags).Error
	if err != nil {
		return err
	}

	// update tag temperature and updated_at
	err = tx.Model(&hole.Tags).Update("temperature", gorm.Expr("temperature + 1")).Error
	return err
}

func (hole *Hole) SetHoleFloor() {
	holeFloorSize := len(hole.Floors)
	if holeFloorSize == 0 {
		return
	}
	hole.HoleFloor.Floors = hole.Floors
	hole.HoleFloor.FirstFloor = hole.Floors[0]
	hole.HoleFloor.LastFloor = hole.Floors[holeFloorSize-1]
	if holeFloorSize > config.Config.HoleFloorSize {
		hole.HoleFloor.Floors = hole.Floors[:holeFloorSize-1]
	}
}

func (hole *Hole) Create(tx *gorm.DB) error {
	// Create hole.Tags, in different sql session
	err := hole.Tags.FindOrCreateTags(tx)
	if err != nil {
		return err
	}

	err = tx.Transaction(func(tx *gorm.DB) error {
		// Create hole
		err = tx.Omit(clause.Associations).Create(hole).Error
		if err != nil {
			return err
		}

		// Create hole_tags association only
		err = tx.Omit("Tags.*", "UpdatedAt").Select("Tags").Save(&hole).Error
		if err != nil {
			return err
		}

		// Update tag temperature
		err = hole.Tags.AddTagTemperature(tx)
		if err != nil {
			return err
		}

		// todo: new Anonyname mapping

		// create floor
		err = hole.Floors[0].Create(tx)
		if err != nil {
			return err
		}

		// create Favorite
		return AddUserFavourite(tx, hole.UserID, hole.ID)
	})
	// transaction commit here
	if err != nil {
		return err
	}

	// set hole.HoleFloor
	hole.SetHoleFloor()

	// store into cache
	return utils.SetCache(hole.IDString(), hole, HoleCacheExpire)
}

func (hole *Hole) AfterCreate(_ *gorm.DB) (err error) {
	hole.HoleID = hole.ID
	return nil
}

func (hole *Hole) AfterFind(_ *gorm.DB) (err error) {
	hole.HoleID = hole.ID
	return nil
}
