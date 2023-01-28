package models

import (
	"fmt"
	"gorm.io/plugin/dbresolver"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Floor struct {
	/// saved fields
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created"`
	UpdatedAt time.Time `json:"time_updated"`

	/// base info

	// content of the floor, no more than 15000
	Content string `json:"content" gorm:"size:15000"`

	// a random username
	Anonyname string `json:"anonyname" gorm:"size:32"`

	// the ranking of this floor in the hole
	Ranking int `json:"ranking" gorm:"default:0;not null;uniqueIndex:idx_hole_ranking,priority:2"`

	// floor_id that it replies to, for dialog mode
	ReplyTo int `json:"reply_to"`

	// like number
	Like int `json:"like" gorm:"not null:default:0"`

	// dislike number
	Dislike int `json:"dislike" gorm:"not null:default:0"`

	// whether the floor is deleted
	Deleted bool `json:"deleted" gorm:"not null;default:false"`

	// the modification times of floor.content
	Modified int `json:"modified" gorm:"not null;default:0"`

	// fold reason
	Fold string `json:"fold_v2"`

	// additional info, like "树洞管理团队"
	SpecialTag string `json:"special_tag"`

	/// association info, should add foreign key

	// the user who wrote it, hidden to user but available to admin
	UserID int `json:"user_id;omitempty"`

	// the hole it belongs to
	HoleID int `json:"hole_id" gorm:"not null;uniqueIndex:idx_hole_ranking,priority:1"`

	// many to many mentions (in different holes)
	Mention Floors `json:"mention" gorm:"many2many:floor_mention"`

	LikedUsers Users `json:"-" gorm:"many2many:floor_like"`

	DislikedUsers Users `json:"-" gorm:"many2many:floor_dislike"`

	// a floor has many history
	History FloorHistorySlice `json:"-"`

	/// dynamically generated fields

	// old version compatibility
	FloorID int `json:"floor_id" gorm:"-:all"`

	// fold reason, for v1
	FoldFrontend []string `json:"fold" gorm:"-:all"`

	// whether the user has liked or disliked the floor
	Liked int8 `json:"-" gorm:"-:all"`

	// whether the user has liked the floor
	LikedFrontend bool `json:"liked" gorm:"-:all"`

	// whether the user has disliked the floor
	DislikedFrontend bool `json:"disliked" gorm:"-:all"`

	// whether the user is the author of the floor
	IsMe bool `json:"is_me" gorm:"-:all"`
}

func (floor *Floor) GetID() int {
	return floor.ID
}

type Floors []*Floor

/******************************
Get and List
*******************************/

func (floor *Floor) Preprocess(c *fiber.Ctx) error {
	return Floors{floor}.Preprocess(c)
}

func (floors Floors) Preprocess(c *fiber.Ctx) error {
	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	// get floors' like
	floorIDs := make([]int, len(floors))
	IDFloorMapping := make(map[int]*Floor)
	for i, floor := range floors {
		if userID == floor.UserID {
			floors[i].IsMe = true
		}
		floorIDs[i] = floor.ID
		IDFloorMapping[floor.ID] = floors[i]
	}

	var floorLikes []FloorLike
	result := DB.
		Clauses(dbresolver.Write).
		Where("floor_id IN (?)", floorIDs).
		Where("user_id = ?", userID).
		Find(&floorLikes)
	if result.Error != nil {
		return err
	}
	for _, floorLike := range floorLikes {
		if floor, ok := IDFloorMapping[floorLike.FloorID]; ok {
			floor.Liked = floorLike.LikeData
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
		for j := range floors[i].Mention {
			floors[i].Mention[j].SetDefaults()
		}
	}
	return nil
}

func (floor *Floor) SetDefaults() {
	if floor.Mention == nil {
		floor.Mention = Floors{}
	}

	floor.FloorID = floor.ID

	if floor.Fold != "" {
		floor.FoldFrontend = []string{floor.Fold}
	} else {
		floor.FoldFrontend = []string{}
	}
}

func (floor *Floor) SetMention(tx *gorm.DB, clear bool) error {
	// find mention IDs
	mentionIDs, err := parseFloorMentions(tx, floor.Content)
	if err != nil {
		return err
	}

	// find mention from floor table
	mention := Floors{}
	if len(mentionIDs) > 0 {
		err := tx.Find(&mention, mentionIDs).Error
		if err != nil {
			return err
		}
	}
	floor.Mention = mention

	// set mention to floor_mention table
	if clear {
		if err = deleteFloorMentions(tx, floor.ID); err != nil {
			return err
		}
	}
	return insertFloorMentions(tx, newFloorMentions(floor.ID, mentionIDs))
}

/******************************
Create
*******************************/

func (floor *Floor) Create(c *fiber.Ctx, db ...*gorm.DB) error {
	var tx *gorm.DB
	if len(db) > 0 {
		tx = db[0]
	} else {
		tx = DB
	}

	// get user
	user, err := GetUserFromAuth(c)
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

		// create floor
		result = tx.Omit("Mention").Create(floor)
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

	// floor set Mention
	err = floor.SetMention(tx, false)
	if err != nil {
		return err
	}

	// update reply and update_at
	result := tx.Exec("UPDATE hole SET reply = reply + 1, updated_at = ? WHERE id = ?", time.Now(), floor.HoleID)
	if result.Error != nil {
		return result.Error
	}

	var messages Messages
	messages = messages.Merge(floor.SendReply(tx))
	messages = messages.Merge(floor.SendMention(tx))
	messages = messages.Merge(floor.SendFavorite(tx))

	err = messages.Send()
	if err != nil {
		utils.Logger.Error("[notification] SendMessage failed: " + err.Error())
		// return err // only for test
	}

	return nil
}

//	Update and Modify

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
	if floor.UserID == userID {
		floor.IsMe = true
	}

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
		if like == 1 {
			floor.LikedFrontend = true
		} else if like == -1 {
			floor.DislikedFrontend = true
		}
		return nil
	})
}

/***************************
Send Notifications
******************/

func (floor *Floor) SendFavorite(tx *gorm.DB) Message {
	// get recipients
	var tmpIDs []int
	result := tx.Raw("SELECT user_id from user_favorites WHERE hole_id = ?", floor.HoleID).Scan(&tmpIDs)
	if result.Error != nil {
		return nil
	}

	// filter my id
	var userIDs []int
	for _, id := range tmpIDs {
		if id != floor.UserID {
			userIDs = append(userIDs, id)
		}
	}

	// return if no recipients
	if len(userIDs) == 0 {
		return nil
	}

	// Construct Message
	message := Message{
		"data":       floor,
		"recipients": userIDs,
		"type":       MessageTypeFavorite,
		"url":        fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	return message
}

func (floor *Floor) SendReply(tx *gorm.DB) Message {
	// get recipients
	userID := 0
	result := tx.Raw("SELECT user_id from hole WHERE id = ?", floor.HoleID).Scan(&userID)
	if result.Error != nil {
		return nil
	}

	// return if no recipients or isMe
	if userID == 0 || userID == floor.UserID {
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

	return message
}

func (floor *Floor) SendMention(_ *gorm.DB) Message {
	// get recipients
	var userIDs []int
	for _, mention := range floor.Mention {
		// not send to me
		if mention.UserID == floor.UserID {
			continue
		}

		userIDs = append(userIDs, mention.UserID)
	}

	// return if no recipients
	if len(userIDs) == 0 {
		return nil
	}

	// construct message
	message := Message{
		"data":       floor,
		"recipients": userIDs,
		"type":       MessageTypeMention,
		"url":        fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	return message
}

func (floor *Floor) SendModify(_ *gorm.DB) error {
	// get recipients
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
