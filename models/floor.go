package models

import (
	"regexp"
	"strconv"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Floor has a tree structure, example:
//	id: 1, reply_to: 0, storey: 1
//		id: 2, reply_to: 1, storey: 2
//	id: 3, reply_to: 0, storey: 3
//		id: 4, reply_to: 3, storey: 4
//			id: 6, reply_to: 4, storey: 5
//		id: 5, reply_to: 3, storey: 6
//	id: 7, reply_to: 0, storey: 7
type Floor struct {
	BaseModel
	HoleID     int     `json:"hole_id"`
	UserID     int     `json:"-"`
	Content    string  `json:"content"`                                // not empty
	Anonyname  string  `json:"anonyname" gorm:"size:32"`               // random username, not empty
	Storey     int     `json:"storey"`                                 // The sequence of floors in a hole
	Path       string  `json:"path" default:"/"`                       // storey path
	ReplyTo    int     `json:"-"`                                      // Floor id that it replies to (must be in the same hole)
	Mention    []Floor `json:"mention" gorm:"many2many:floor_mention"` // Many to many mentions (in different holes)
	Like       int     `json:"like"`                                   // like - dislike
	Liked      int8    `json:"liked" gorm:"-:all"`                     // whether the user has liked or disliked the floor, dynamically generated
	IsMe       bool    `json:"is_me" gorm:"-:all"`                     // whether the user is the author of the floor, dynamically generated
	Deleted    bool    `json:"deleted"`                                // whether the floor is deleted
	Fold       string  `json:"fold"`                                   // fold reason
	SpecialTag string  `json:"special_tag"`                            // Additional info
}

type Floors []Floor

type AnonynameMapping struct {
	HoleID    int    `json:"hole_id" gorm:"primarykey"`
	UserID    int    `json:"user_id" gorm:"primarykey"`
	Anonyname string `json:"anonyname" gorm:"index;size:32"`
}

type FloorLike struct {
	FloorID  int  `json:"floor_id" gorm:"primarykey"`
	UserID   int  `json:"user_id" gorm:"primarykey"`
	LikeData int8 `json:"like_data"`
}

//goland:noinspection GoNameStartsWithPackageName
type FloorHistory struct {
	BaseModel
	Content string `json:"content"`
	Reason  string `json:"reason"`
	FloorID int    `json:"floor_id"`
	UserID  int    `json:"user_id"` // The one who modified the floor
}

func (floor *Floor) Preprocess() error {
	return floor.LoadMention()
}

func (floors Floors) Preprocess() error {
	for i := 0; i < len(floors); i++ {
		err := floors[i].Preprocess()
		if err != nil {
			return err
		}
	}
	return nil
}

func (floor *Floor) LoadMention() error {
	var Mention Floors
	err := DB.Model(floor).Association("Mention").Find(&Mention)
	if err != nil {
		return err
	}
	if Mention != nil {
		floor.Mention = Mention
	}
	return nil
}

func (floor Floor) MakeQuerySet(
	limit int, offset int,
	holeID int, orderBy string,
	ifDesc bool) (tx *gorm.DB) {
	return DB.
		Limit(limit).Offset(offset).
		Where("hole_id = ?", holeID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: orderBy}, Desc: ifDesc})
}

var reHole = regexp.MustCompile(`[^#]#(\d+)`)
var reFloor = regexp.MustCompile(`##(\d+)`)

func (floor *Floor) FindMention(tx *gorm.DB) error {
	holeIDsText := reHole.FindAllStringSubmatch(" "+floor.Content, -1)
	holeIds, err := utils.ReText2IntArray(holeIDsText)
	if err != nil {
		return err
	}

	var floorIDs []int
	if len(holeIds) != 0 {
		err := tx.
			Raw("SELECT MIN(id) FROM floor WHERE hole_id IN ? GROUP BY hole_id", holeIds).
			Scan(&floorIDs).Error
		if err != nil {
			return err
		}
	}

	floorIDsText := reFloor.FindAllStringSubmatch(" "+floor.Content, -1)
	floorIDs2, err := utils.ReText2IntArray(floorIDsText)
	if err != nil {
		return err
	}

	floorIDs = append(floorIDs, floorIDs2...)
	var floors []Floor
	if len(floorIDs) != 0 {
		err := tx.Find(&floors, floorIDs).Error
		if err != nil {
			return err
		}
	}
	floor.Mention = floors

	return nil
}

func (floor *Floor) Create(c *fiber.Ctx, db ...*gorm.DB) error {
	var tx *gorm.DB
	if len(db) > 0 {
		tx = db[0]
	} else {
		tx = DB
	}

	userID, err := GetUserID(c)
	if err != nil {
		return err
	}
	floor.UserID = userID
	floor.IsMe = true

	err = tx.Transaction(func(tx *gorm.DB) error {
		// get anonymous name
		var mapping AnonynameMapping

		result := tx.
			Where("hole_id = ?", floor.HoleID).
			Where("user_id = ?", userID).
			Take(&mapping)

		if result.Error != nil {
			// no mapping exists, generate anonyname
			var names []string
			result = tx.Clauses(clause.Locking{
				Strength: "UPDATE",
			}).Raw(`
				SELECT anonyname FROM anonyname_mapping 
				WHERE hole_id = ? 
				ORDER BY anonyname`, floor.HoleID,
			).Scan(&names)
			if result.Error != nil {
				return result.Error
			}

			floor.Anonyname = utils.GenerateName(names)
			result = tx.Create(&AnonynameMapping{
				UserID:    userID,
				HoleID:    floor.HoleID,
				Anonyname: floor.Anonyname,
			})
			if result.Error != nil {
				return result.Error
			}
		} else {
			floor.Anonyname = mapping.Anonyname
		}

		// set storey and path
		if floor.ReplyTo == 0 {
			var count int64
			result := tx.Clauses(clause.Locking{
				Strength: "UPDATE",
			}).Model(&Floor{}).Where("hole_id = ?", floor.HoleID).
				Count(&count)
			if result.Error != nil {
				return err
			}
			floor.Storey = int(count) + 1
			floor.Path = "/"
		} else {
			var storey int
			result := tx.Clauses(clause.Locking{
				Strength: "UPDATE",
			}).Raw(`SELECT storey FROM floor 
					WHERE hole_id = ? AND path LIKE '%/?/%' 
					ORDER BY storey DESC LIMIT 1`,
				floor.HoleID, floor.ReplyTo).
				Scan(&storey)
			if result.Error != nil {
				return err
			}
			result = tx.
				Exec(`UPDATE floor SET storey = storey + 1
					WHERE hole_id = ? AND storey > ?`,
					floor.HoleID, storey)
			if result.Error != nil {
				return err
			}
			floor.Storey = storey + 1
			var replyPath string
			result = tx.
				Raw(`SELECT path FROM floor WHERE ID = ?`,
					floor.ReplyTo).
				Scan(&replyPath)
			if result.Error != nil {
				return err
			}
			floor.Path = replyPath + strconv.Itoa(floor.ReplyTo) + "/"
		}

		// create floor
		result = tx.Create(floor)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (floor *Floor) BeforeCreate(tx *gorm.DB) (err error) {
	// find mention
	err = floor.FindMention(tx)
	if err != nil {
		return err
	}
	return nil
}

func (floor *Floor) AfterCreate(tx *gorm.DB) (err error) {
	result := tx.Exec("UPDATE hole SET reply = reply + 1 WHERE id = ?", floor.HoleID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (floor *Floor) Backup(c *fiber.Ctx, reason string) error {
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	history := FloorHistory{
		Content: floor.Content,
		Reason:  reason,
		FloorID: floor.ID,
		UserID:  userID,
	}
	return DB.Create(&history).Error
}

func (floor *Floor) ModifyLike(c *fiber.Ctx, likeOption int8) error {
	// validate like option
	if likeOption > 1 || likeOption < -1 {
		return utils.BadRequest("like option must be -1, 0 or 1")
	}

	// get userID
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}
	floor.IsMe = true

	return DB.Transaction(func(tx *gorm.DB) error {

		result := tx.Exec("DELETE FROM floor_like WHERE floor_id = ? AND user_id = ?", floor.ID, userID)
		if result.Error != nil {
			return result.Error
		}

		if likeOption != 0 {
			result = tx.Create(&FloorLike{
				FloorID:  floor.ID,
				UserID:   userID,
				LikeData: likeOption,
			})
			if result.Error != nil {
				return err
			}
		}

		var like int
		result = tx.Raw(`
			SELECT IFNULL(SUM(like_data), 0)
			FROM floor_like 
			WHERE floor_id = ?`,
			floor.ID,
		).Scan(&like)
		if result.Error != nil {
			return result.Error
		}

		floor.Like = like
		floor.Liked = likeOption
		return nil
	})
}

func (floor *Floor) LoadDyField(c *fiber.Ctx) error {
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}
	if userID == floor.UserID {
		floor.IsMe = true
	}

	var floorLike FloorLike
	result := DB.
		Where("floor_id = ?", floor.ID).
		Where("user_id = ?", userID).
		Take(&floorLike)
	if result.Error == nil {
		floor.Liked = floorLike.LikeData
	}
	return nil
}

func (floors Floors) LoadDyField(c *fiber.Ctx) error {
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}
	var floorIDs []int
	IDFloorMapping := make(map[int]*Floor)
	for i, v := range floors {
		if userID == v.UserID {
			floors[i].IsMe = true
		}
		floorIDs = append(floorIDs, v.ID)
		IDFloorMapping[v.ID] = &floors[i]
	}

	var floorLikes []FloorLike
	result := DB.
		Where("floor_id IN ?", &floorIDs).
		Where("user_id = ?", userID).
		Find(&floorLikes)
	if result.Error != nil {
		return err
	}
	for _, v := range floorLikes {
		if floor, ok := IDFloorMapping[v.FloorID]; ok {
			floor.Liked = v.LikeData
		}
	}
	return nil
}
