package models

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"treehole_next/config"
	"treehole_next/utils"

)

type HoleV2 struct {
	BaseHole

	PrefetchFloors Floors `json:"prefetch_floors" gorm:"-:all"`
}

func (hole *HoleV2) SetHoleFloor() {
	if len(hole.Floors) != 0 {
		holeFloorSize := len(hole.Floors)

		if holeFloorSize <= config.Config.HoleFloorSize {
			hole.PrefetchFloors = hole.Floors
		} else {
			hole.PrefetchFloors = hole.Floors[0 : holeFloorSize - 1]
		}
	} else if len(hole.PrefetchFloors) != 0 {
		hole.Floors = hole.PrefetchFloors
	}
}

type HolesV2 []*HoleV2

func (holes HolesV2) loadFloors() error {
	if len(holes) == 0 {
		return nil
	}
	holeIDs := utils.Models2IDSlice(holes)

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
				DB.Model(&BaseHole{}).Select("id", "reply").Where("id in ?", holeIDs),
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
	index := slices.IndexFunc(holes, func(hole *HoleV2) bool {
		return hole.ID == floors[0].HoleID
	})
	for _, floor := range floors {
		if floor.HoleID != holes[index].ID {
			holes[index].Floors = floors[left:right]
			left = right
			index = slices.IndexFunc(holes, func(hole *HoleV2) bool {
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

func (hole *HoleV2) Create(tx *gorm.DB, user *User, tagNames []string, c *fiber.Ctx) (err error) {
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
		firstFloor.SendSensitive(tx)
		// firstFloor.Content = ""
	}

	hole.HoleHook()

	// store into cache
	return utils.SetCache(hole.CacheName(), hole, HoleCacheExpire)
}

func (hole *HoleV2) HoleHook() {
	if hole == nil {
		return
	}
	notifyMessage := fmt.Sprintf("#%d\n", hole.ID)

	if hole.DivisionID == 4 {
		go utils.NotifyQQ(&utils.BotMessage{
			MessageType: utils.MessageTypePrivate,
			UserID:      config.Config.QQBotUserID,
			Message:     notifyMessage,
		})
		go utils.NotifyFeishu(&utils.FeishuMessage{
			MsgType: "text",
			Content: notifyMessage,
		})
	}

	tagToGroup := map[string]*int64{
		"@物理大神": config.Config.QQBotPhysicsGroupID,
		"@码上辅导": config.Config.QQBotCodingGroupID,
	}

	for _, tag := range hole.Tags {
		if tag != nil {
			if groupID, ok := tagToGroup[tag.Name]; ok {
				go utils.NotifyQQ(&utils.BotMessage{
					MessageType: utils.MessageTypeGroup,
					GroupID:     groupID,
					Message:     notifyMessage,
				})
			}
		}
	}
}