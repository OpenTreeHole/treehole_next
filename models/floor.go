package models

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Floor has a tree structure, example:
//
//	id: 1, reply_to: 0, storey: 1
//		id: 2, reply_to: 1, storey: 2
//	id: 3, reply_to: 0, storey: 3
//		id: 4, reply_to: 3, storey: 4
//			id: 6, reply_to: 4, storey: 5
//		id: 5, reply_to: 3, storey: 6
//	id: 7, reply_to: 0, storey: 7
type Floor struct {
	BaseModel
	FloorID          int      `json:"floor_id" gorm:"-:all"`
	HoleID           int      `json:"hole_id"`                                // the hole it belongs to
	UserID           int      `json:"-"`                                      // the user who wrote it, hidden to user but available to admin
	Content          string   `json:"content"`                                // content of the floor
	Anonyname        string   `json:"anonyname" gorm:"size:32"`               // a random username
	Storey           int      `json:"storey"`                                 // the sequence of floors in a hole
	Path             string   `json:"path" gorm:"default:/"`                  // storey path e.g. /1/2/3/
	ReplyTo          int      `json:"-" gorm:"-:all"`                         // Floor id that it replies to (must be in the same hole)
	Mention          []Floor  `json:"mention" gorm:"many2many:floor_mention"` // many to many mentions (in different holes)
	Like             int      `json:"like"`                                   // like number - dislike number
	Liked            int8     `json:"-" gorm:"-:all"`                         // whether the user has liked or disliked the floor, dynamically generated
	LikedFrontend    bool     `json:"liked" gorm:"-:all"`                     // whether the user has liked the floor, dynamically generated
	DislikedFrontend bool     `json:"disliked" gorm:"-:all"`                  // whether the user has disliked the floor, dynamically generated
	IsMe             bool     `json:"is_me" gorm:"-:all"`                     // whether the user is the author of the floor, dynamically generated
	Deleted          bool     `json:"deleted"`                                // whether the floor is deleted
	Fold             string   `json:"fold_v2"`                                // fold reason
	FoldFrontend     []string `json:"fold" gorm:"-:all"`                      // fold reason, for v1
	SpecialTag       string   `json:"special_tag"`                            // additional info, like "树洞管理团队"
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

func (floor *Floor) Preprocess(c *fiber.Ctx) error {
	floors := Floors{*floor}

	err := floors.Preprocess(c)
	if err != nil {
		return err
	}

	*floor = floors[0]

	return nil
}

func (floors Floors) Preprocess(c *fiber.Ctx) error {
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// get floors' like
	floorIDs := make([]int, len(floors))
	IDFloorMapping := make(map[int]*Floor)
	for i, v := range floors {
		if userID == v.UserID {
			floors[i].IsMe = true
		}
		floorIDs[i] = v.ID
		IDFloorMapping[v.ID] = &floors[i]
	}

	var floorLikes []FloorLike
	result := DB.
		Where("floor_id IN (?)", floorIDs).
		Where("user_id = ?", userID).
		Find(&floorLikes)
	if result.Error != nil {
		return err
	}
	for _, v := range floorLikes {
		if floor, ok := IDFloorMapping[v.FloorID]; ok {
			floor.Liked = v.LikeData
			switch floor.Liked {
			case 1:
				floor.LikedFrontend = true
			case -1:
				floor.DislikedFrontend = true
			}
		}
	}

	// set some default values
	for i := range floors {
		floors[i].SetDefaults()
	}
	return nil
}

func (floor *Floor) SetDefaults() {
	if floor.Mention == nil {
		floor.Mention = []Floor{}
	}

	floor.FloorID = floor.ID

	if floor.Fold != "" {
		floor.FoldFrontend = []string{floor.Fold}
	} else {
		floor.FoldFrontend = []string{}
	}
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
	floors := make([]Floor, 0)
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

	// get user
	var user User
	err := user.GetUser(c)
	if err != nil {
		return err
	}
	floor.UserID = user.ID
	floor.IsMe = true

	return tx.Transaction(func(tx *gorm.DB) error {
		// permission
		var hole Hole
		tx.Select("division_id").First(&hole, floor.HoleID)
		if user.BanDivision[hole.DivisionID] ||
			floor.SpecialTag != "" && !perm.CheckPermission(user, perm.Admin|perm.Operator) {
			return utils.Forbidden()
		}

		// get anonymous name
		var mapping AnonynameMapping
		result := tx.
			Where("hole_id = ?", floor.HoleID).
			Where("user_id = ?", floor.UserID).
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
				UserID:    floor.UserID,
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
			result = tx.Clauses(clause.Locking{
				Strength: "UPDATE",
			}).Model(&Floor{}).Where("hole_id = ?", floor.HoleID).
				Count(&count)
			if result.Error != nil {
				return result.Error
			}
			floor.Storey = int(count) + 1
			floor.Path = "/"
		} else {
			var storey int
			result = tx.Clauses(clause.Locking{
				Strength: "UPDATE",
			}).Raw(`SELECT storey FROM floor 
					WHERE hole_id = ? AND path LIKE '%/?/%' 
					ORDER BY storey DESC LIMIT 1`,
				floor.HoleID, floor.ReplyTo).
				Scan(&storey)
			if result.Error != nil {
				return result.Error
			}
			result = tx.
				Exec(`UPDATE floor SET storey = storey + 1
					WHERE hole_id = ? AND storey > ?`,
					floor.HoleID, storey)
			if result.Error != nil {
				return result.Error
			}
			floor.Storey = storey + 1
			var replyPath string
			result = tx.
				Raw(`SELECT path FROM floor WHERE ID = ?`,
					floor.ReplyTo).
				Scan(&replyPath)
			if result.Error != nil {
				return result.Error
			}
			floor.Path = replyPath + strconv.Itoa(floor.ReplyTo) + "/"
		}
		// find mention
		err = floor.FindMention(tx)
		if err != nil {
			return err
		}

		// create floor
		result = tx.Create(floor)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&Hole{}).Where("id = ?", floor.HoleID).Update("updated_at", time.Now())
		if result.Error != nil {
			return result.Error
		}

		if hole.Reply < config.Config.HoleFloorSize {
			return utils.DeleteCache(fmt.Sprintf("hole_%d", floor.HoleID))
		}
		return nil
	})
}

func (floor *Floor) AfterCreate(tx *gorm.DB) (err error) {

	result := tx.Exec("UPDATE hole SET reply = reply + 1 WHERE id = ?", floor.HoleID)
	if result.Error != nil {
		return result.Error
	}

	err = floor.SendFavorite(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendFavorite failed: " + err.Error())
		// return err // only for test
	}

	err = floor.SendReply(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendReply failed: " + err.Error())
		// return err // only for test
	}

	err = floor.SendMention(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendMention failed: " + err.Error())
		// return err // only for test
	}

	return nil
}

func (floor *Floor) AfterUpdate(tx *gorm.DB) (err error) {

	err = floor.SendModify(tx)
	if err != nil {
		utils.Logger.Error("[notification] SendModify failed: " + err.Error())
		// return err // only for test
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

func (floor *Floor) SendFavorite(tx *gorm.DB) error {
	// get recipents
	var userIDs []int
	result := tx.Raw("SELECT user_id from user_favorites WHERE hole_id = ?", floor.HoleID).Scan(&userIDs)
	if result.Error != nil {
		return result.Error
	}

	// return if no recipents
	if userIDs == nil {
		return nil
	}

	// Construct Message
	message := Message{
		"data":       floor,
		"recipients": userIDs,
		"type":       MessageTypeFavorite,
		"url":        fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	// Send
	err := message.Send()
	if err != nil {
		return err
	}

	return nil
}

func (floor *Floor) SendReply(tx *gorm.DB) error {
	// get recipents
	userID := 0
	result := tx.Raw("SELECT user_id from hole WHERE id = ?", floor.HoleID).Scan(&userID)
	if result.Error != nil {
		return result.Error
	}

	// return if no recipents
	if userID == 0 {
		return nil
	}

	userIDs := []int{userID}

	// construct message
	message := Message{
		"data":       floor,
		"recipients": userIDs,
		"type":       MessageTypeReply,
		"url":        fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	// send
	err := message.Send()
	if err != nil {
		return err
	}

	return nil
}

func (floor *Floor) SendMention(tx *gorm.DB) error {
	// get recipents
	var userIDs []int
	for _, mention := range floor.Mention {
		userIDs = append(userIDs, mention.UserID)
	}

	// return if no recipents
	if userIDs == nil {
		return nil
	}

	// construct message
	message := Message{
		"data":       floor,
		"recipients": userIDs,
		"type":       MessageTypeMention,
		"url":        fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	// send
	err := message.Send()
	if err != nil {
		return err
	}

	return nil
}

func (floor *Floor) SendModify(tx *gorm.DB) error {
	// get recipents
	userIDs := []int{floor.UserID}

	// construct message
	message := Message{
		"data":       floor,
		"recipients": userIDs,
		"type":       MessageTypeModify,
		"url":        fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	// send
	err := message.Send()
	if err != nil {
		return err
	}

	return nil
}
