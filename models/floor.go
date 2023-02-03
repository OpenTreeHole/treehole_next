package models

import (
	"fmt"
	"time"
	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/plugin/dbresolver"

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

	// floor_id that it replies to, for dialog mode, in the same hole
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

	// the user who wrote it
	UserID int `json:"-" gorm:"not null"`

	// the hole it belongs to
	HoleID int `json:"hole_id" gorm:"not null;uniqueIndex:idx_hole_ranking,priority:1"`

	// many to many mentions
	Mention Floors `json:"mention" gorm:"many2many:floor_mention"`

	LikedUsers Users `json:"-" gorm:"many2many:floor_like"`

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
		floors[i].IsMe = userID == floor.UserID
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
	for _, floor := range floors {
		floor.SetDefaults()
	}
	return nil
}

func (floor *Floor) SetDefaults() {
	if floor.Mention == nil {
		floor.Mention = Floors{}
	} else if len(floor.Mention) > 0 {
		for _, mentionFloor := range floor.Mention {
			mentionFloor.SetDefaults()
		}
	}

	if floor.Fold != "" {
		floor.FoldFrontend = []string{floor.Fold}
	} else {
		floor.FoldFrontend = []string{}
	}
}

/******************************
Create
*******************************/

func (floor *Floor) Create(tx *gorm.DB) (err error) {
	// load floor mention, in another session
	floor.Mention, err = LoadFloorMentions(DB, floor.Content)
	if err != nil {
		return err
	}
	var hole Hole

	err = tx.Transaction(func(tx *gorm.DB) error {
		// get anonymous name
		floor.Anonyname, err = FindOrGenerateAnonyname(tx, floor.HoleID, floor.UserID)
		if err != nil {
			return err
		}

		// get and lock hole for updating reply
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&hole, floor.HoleID).Error
		if err != nil {
			return err
		}

		hole.Reply++
		floor.Ranking = hole.Reply

		// create floor, set floor_mention association in AfterCreate hook
		err = tx.Omit(clause.Associations).Create(&floor).Error
		if err != nil {
			return err
		}

		// update hole reply and update_at
		return tx.Omit(clause.Associations).Save(&hole).Error
	})

	if err != nil {
		return err
	}

	floor.SetDefaults()

	// delete cache
	return utils.DeleteCache(hole.CacheName())
}

func (floor *Floor) AfterFind(_ *gorm.DB) (err error) {
	floor.FloorID = floor.ID
	return nil
}

func (floor *Floor) AfterCreate(tx *gorm.DB) (err error) {
	floor.FloorID = floor.ID

	// create floor mention association
	if len(floor.Mention) > 0 {
		err = tx.Omit("Mention.*", "UpdatedAt").Select("Mention").Save(&floor).Error
		if err != nil {
			return err
		}
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

// ModifyLike do in transaction only
func (floor *Floor) ModifyLike(tx *gorm.DB, userID int, likeOption int8) (err error) {
	if userID == floor.UserID {
		floor.IsMe = true
	}
	floorLike := &FloorLike{
		FloorID: floor.ID,
		UserID:  userID,
	}
	if likeOption == 0 {
		err = tx.Delete(&floorLike).Error
		if err != nil {
			return err
		}
	} else {
		floorLike.LikeData = likeOption
		err = tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&floorLike).Error
		if err != nil {
			return err
		}
	}

	var like, dislike int64
	err = tx.Model(&FloorLike{}).Where("floor_id = ? and like_data = ?", floor.ID, 1).Count(&like).Error
	if err != nil {
		return err
	}
	err = tx.Model(&FloorLike{}).Where("floor_id = ? and like_data = ?", floor.ID, 1).Count(&dislike).Error
	if err != nil {
		return err
	}

	floor.Like = int(like)
	floor.Dislike = int(dislike)
	floor.Liked = likeOption
	if likeOption == 1 {
		floor.LikedFrontend = true
	} else if likeOption == -1 {
		floor.DislikedFrontend = true
	}
	return nil
}

/***************************
Send Notifications
******************/

func (floor *Floor) SendFavorite(tx *gorm.DB) Notification {
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

	// Construct Notification
	message := Notification{
		"data":        floor,
		"recipients":  userIDs,
		"description": floor.Content,
		"title":       "您收藏的树洞有新回复",
		"type":        MessageTypeFavorite,
		"url":         fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	return message
}

func (floor *Floor) SendReply(tx *gorm.DB) Notification {
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
	message := Notification{
		"data":        floor,
		"recipients":  userIDs,
		"description": floor.Content,
		"title":       "您的帖子被回复了",
		"type":        MessageTypeReply,
		"url":         fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	return message
}

func (floor *Floor) SendMention(_ *gorm.DB) Notification {
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
	message := Notification{
		"data":        floor,
		"recipients":  userIDs,
		"description": floor.Content,
		"title":       "您的帖子被引用了",
		"type":        MessageTypeMention,
		"url":         fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	return message
}

func (floor *Floor) SendModify(_ *gorm.DB) error {
	// get recipients
	userIDs := []int{floor.UserID}

	// construct message
	message := Notification{
		"data":        floor,
		"recipients":  userIDs,
		"description": floor.Content,
		"title":       "您的帖子被修改了",
		"type":        MessageTypeModify,
		"url":         fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	// send
	_, err := message.Send()
	if err != nil {
		return err
	}

	return nil
}
