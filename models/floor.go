package models

import (
	"fmt"
	"time"
	"treehole_next/utils/sensitive"

	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"

	"treehole_next/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/plugin/dbresolver"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Floor struct {
	/// saved fields
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"time_created" gorm:"index:,sort:desc"`
	UpdatedAt time.Time `json:"time_updated" gorm:"index:,sort:desc;index:idx_floor_actual_sensitive_updated_at,priority:3,sort:desc"`

	/// base info

	// content of the floor, no more than 15000, should be sensitive checked, no more than 10000 in frontend
	Content string `json:"content" gorm:"not null;size:15000"`

	// a random username
	Anonyname string `json:"anonyname" gorm:"not null;size:32"`

	// the ranking of this floor in the hole
	Ranking int `json:"ranking" gorm:"default:0;not null;uniqueIndex:idx_hole_ranking,priority:2"`

	// floor_id that it replies to, for dialog mode, in the same hole
	ReplyTo int `json:"reply_to" gorm:"not null;default:0"`

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

	// auto sensitive check
	IsSensitive bool `json:"is_sensitive" gorm:"index:idx_floor_actual_sensitive_updated_at,priority:1"`

	// manual sensitive check
	IsActualSensitive *bool `json:"is_actual_sensitive" gorm:"index:idx_floor_actual_sensitive_updated_at,priority:2"`

	// auto sensitive check detail
	SensitiveDetail string `json:"sensitive_detail,omitempty"`

	/// association info, should add foreign key

	// the user who wrote it
	UserID int `json:"-" gorm:"not null"`

	// the hole it belongs to
	HoleID int `json:"hole_id" gorm:"not null;uniqueIndex:idx_hole_ranking,priority:1"`

	// many to many mentions
	Mention Floors `json:"mention" gorm:"many2many:floor_mention;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	LikedUsers Users `json:"-" gorm:"many2many:floor_like;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// a floor has many history
	History FloorHistorySlice `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

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

func MakeFloorQuerySet(_ *fiber.Ctx) (*gorm.DB, error) {
	return DB.Preload("Mention"), nil
	//user, err := GetUser(c)
	//if err != nil {
	//	return nil, err
	//}
	//if user.IsAdmin {
	//	return DB.Preload("Mention"), nil
	//} else {
	//	userID, err := common.GetUserID(c)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return DB.Where("(is_sensitive = 0 AND is_actual_sensitive IS NULL) OR is_actual_sensitive = 0 OR user_id = ?", userID).
	//		Preload("Mention", "(is_sensitive = 0 AND is_actual_sensitive IS NULL) OR is_actual_sensitive = 0 OR user_id = ?", userID), nil
	//}
}

func (floors Floors) MakeQuerySet(holeID *int, offset, size *int, c *fiber.Ctx) (*gorm.DB, error) {
	querySet, err := MakeFloorQuerySet(c)
	if err != nil {
		return nil, err
	}
	if holeID != nil {
		querySet = querySet.Where("hole_id = ?", holeID)
	}

	if offset != nil {
		querySet = querySet.Offset(*offset)
	}
	if size != nil {
		querySet = querySet.Limit(*size)
	}
	return querySet, nil
}

func (floors Floors) loadFloorLikes(c *fiber.Ctx) (err error) {
	userID, err := common.GetUserID(c)
	if err != nil {
		return
	}

	floorIDs := make([]int, len(floors))
	IDFloorMapping := make(map[int]*Floor)
	for i, floor := range floors {
		floorIDs[i] = floor.ID
		IDFloorMapping[floor.ID] = floors[i]
	}

	var floorLikes []FloorLike
	err = DB.Clauses(dbresolver.Write).
		Where("floor_id IN (?)", floorIDs).
		Where("user_id = ?", userID).
		Find(&floorLikes).Error
	if err != nil {
		return
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
	return
}

func (floors Floors) Preprocess(c *fiber.Ctx) (err error) {
	userID, err := common.GetUserID(c)
	if err != nil {
		return
	}

	// get floors' like
	err = floors.loadFloorLikes(c)
	if err != nil {
		return
	}

	// set floors IsMe
	for _, floor := range floors {
		floor.IsMe = userID == floor.UserID
	}

	// set some default values
	for _, floor := range floors {
		err = floor.SetDefaults(c)
		if err != nil {
			return
		}
	}
	return
}

func (floor *Floor) SetDefaults(c *fiber.Ctx) (err error) {
	floor.FloorID = floor.ID
	user, err := GetUser(c)
	if err != nil {
		return
	}

	floor.Anonyname = utils.GetFuzzName(floor.Anonyname)
	if floor.Sensitive() {
		if user.IsAdmin {
			floor.SpecialTag = "sensitive"
		}
		if !floor.Deleted {
			if floor.IsActualSensitive != nil && *floor.IsActualSensitive {
				// deprecated, deleted already
				floor.Content = "该内容因违反社区规范被删除"
				floor.Deleted = true
			} else {
				floor.Content = "该内容正在审核中"
			}
			floor.FoldFrontend = []string{floor.Content}
			floor.Fold = floor.Content
		}
	}
	if !user.IsAdmin {
		floor.SensitiveDetail = ""
	}

	if floor.Mention == nil {
		floor.Mention = Floors{}
	} else if len(floor.Mention) > 0 {
		for _, mentionFloor := range floor.Mention {
			err = mentionFloor.SetDefaults(c)
			if err != nil {
				return
			}
		}
	}

	if floor.Fold != "" {
		floor.FoldFrontend = []string{floor.Fold}
	} else {
		floor.FoldFrontend = []string{}
	}

	// 直接清空内容而不是替换
	//if !floor.Deleted &&
	//	floor.IsSensitive == true {
	//	var alterContent string
	//	if floor.IsActualSensitive == nil {
	//		alterContent = "该内容被猫猫吃掉了"
	//	} else if *floor.IsActualSensitive == true {
	//		alterContent = "该内容因违反社区规范被删除"
	//	} else {
	//		alterContent = ""
	//	}
	//
	//	if alterContent != "" {
	//		floor.Content = alterContent
	//		floor.FoldFrontend = []string{alterContent}
	//	}
	//}
	return
}

/******************************
Create
*******************************/

func (floor *Floor) Create(tx *gorm.DB, hole *Hole, c *fiber.Ctx) (err error) {
	// sensitive check
	sensitiveCheckResp, err := sensitive.CheckSensitive(sensitive.ParamsForCheck{
		Content:  floor.Content,
		Id:       time.Now().UnixNano(),
		TypeName: sensitive.TypeFloor,
	})
	if err != nil {
		return
	}
	floor.IsSensitive = !sensitiveCheckResp.Pass
	floor.SensitiveDetail = sensitiveCheckResp.Detail

	// load floor mention, in another session
	floor.Mention, err = LoadFloorMentions(DB, floor.Content)
	if err != nil {
		return
	}

	err = tx.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
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
		return tx.Model(&hole).
			Omit(clause.Associations).
			Select("Reply").
			Updates(&hole).Error
	})

	if err != nil {
		return err
	}

	err = floor.SetDefaults(c)
	if err != nil {
		return err
	}

	if !floor.Sensitive() {
		// Send Notification
		var messages Notifications
		messages = messages.Merge(floor.SendReply(tx))
		messages = messages.Merge(floor.SendMention(tx))
		messages = messages.Merge(floor.SendSubscription(tx))

		err = messages.Send()
		if err != nil {
			log.Err(err).Str("model", "Notification").Msg("SendNotification failed")
			// return err // only for test
		}
	}

	if !hole.Hidden && !floor.Sensitive() {
		// insert into Elasticsearch
		go FloorIndex(FloorModel{
			ID:        floor.ID,
			UpdatedAt: time.Now(),
			Content:   floor.Content,
		})
	} else {
		go FloorDelete(floor.ID)
	}

	// delete cache
	return utils.DeleteCache(hole.CacheName())
}

func (floor *Floor) Sensitive() bool {
	if floor == nil {
		return false
	}
	if floor.IsActualSensitive != nil {
		return *floor.IsActualSensitive
	}
	return floor.IsSensitive
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

// Backup Update and Modify
func (floor *Floor) Backup(tx *gorm.DB, userID int, reason string) error {
	history := FloorHistory{
		Content:           floor.Content,
		Reason:            reason,
		FloorID:           floor.ID,
		UserID:            userID,
		IsSensitive:       floor.IsSensitive,
		IsActualSensitive: floor.IsActualSensitive,
		SensitiveDetail:   floor.SensitiveDetail,
	}
	return tx.Create(&history).Error
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
	err = tx.Model(&FloorLike{}).Where("floor_id = ? and like_data = ?", floor.ID, -1).Count(&dislike).Error
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

func (floor *Floor) SendSubscription(tx *gorm.DB) Notification {
	// get recipients
	var tmpIDs []int
	result := tx.Raw("SELECT user_id from user_subscription WHERE hole_id = ?", floor.HoleID).Scan(&tmpIDs)
	if result.Error != nil {
		tmpIDs = []int{}
	}

	// filter my id
	var userIDs []int
	for _, id := range tmpIDs {
		if id != floor.UserID {
			userIDs = append(userIDs, id)
		}
	}

	// Construct Notification
	message := Notification{
		Data:        floor,
		Recipients:  userIDs,
		Description: floor.Content,
		Title:       "您关注的树洞有新回复",
		Type:        MessageTypeFavorite,
		URL:         fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	return message
}

func (floor *Floor) SendReply(tx *gorm.DB) Notification {
	// get recipients
	userID := 0
	result := tx.Raw("SELECT user_id from hole WHERE id = ?", floor.HoleID).Scan(&userID)
	if result.Error != nil {
		userID = 0
	}

	// return if no recipients or isMe
	var userIDs []int
	if userID != 0 && userID != floor.UserID {
		userIDs = []int{userID}
	}

	// construct message
	message := Notification{
		Data:        floor,
		Recipients:  userIDs,
		Description: floor.Content,
		Title:       "您的内容有新回复",
		Type:        MessageTypeReply,
		URL:         fmt.Sprintf("/api/floors/%d", floor.ID),
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

	// construct message
	message := Notification{
		Data:        floor,
		Recipients:  userIDs,
		Description: floor.Content,
		Title:       "您的内容被引用了",
		Type:        MessageTypeMention,
		URL:         fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	return message
}

func (floor *Floor) SendModify(_ *gorm.DB) error {
	// get recipients
	userIDs := []int{floor.UserID}

	// construct message
	message := Notification{
		Data:        floor,
		Recipients:  userIDs,
		Description: floor.Content,
		Title:       "您的内容被管理员修改了",
		Type:        MessageTypeModify,
		URL:         fmt.Sprintf("/api/floors/%d", floor.ID),
	}

	// send
	_, err := message.Send()
	if err != nil {
		return err
	}

	return nil
}
