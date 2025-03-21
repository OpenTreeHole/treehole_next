package models

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"treehole_next/config"

	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type User struct {
	/// base info
	ID int `json:"id" gorm:"primaryKey"`

	Config UserConfig `json:"config" gorm:"serializer:json;not null;default:\"{}\""`

	BanDivision map[int]*time.Time `json:"-" gorm:"serializer:json;not null;default:\"{}\""`

	OffenceCount int `json:"-" gorm:"not null;default:0"`

	BanReport *time.Time `json:"-" gorm:"serializer:json"`

	BanReportCount int `json:"-" gorm:"not null;default:0"`

	DefaultSpecialTag string `json:"default_special_tag" gorm:"size:32"`

	SpecialTags []string `json:"special_tags" gorm:"serializer:json;not null;default:\"[]\""`

	FavoriteGroupCount int `json:"favorite_group_count" gorm:"not null;default:0"`

	/// association fields, should add foreign key

	// holes owned by the user
	UserHoles Holes `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// floors owned by the user
	UserFloors Floors `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// reports made by the user; a user has many report
	UserReports Reports `json:"-"`

	// floors liked by the user
	UserLikedFloors Floors `json:"-" gorm:"many2many:floor_like;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// floor history made by the user
	UserFloorHistory FloorHistorySlice `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// user punishments on division
	UserPunishments Punishments `json:"-"`

	// punishments made by this user
	UserMakePunishments Punishments `json:"-" gorm:"foreignKey:MadeBy"`

	// user punishments on report
	UserReportPunishments ReportPunishments `json:"-"`

	// report punishments made by this user
	UserMakeReportPunishments ReportPunishments `json:"-" gorm:"foreignKey:MadeBy"`

	UserSubscription Holes `json:"-" gorm:"many2many:user_subscription;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	/// dynamically generated field

	UserID int `json:"user_id" gorm:"-:all"`

	Permission struct {
		// 管理员权限到期时间
		Admin time.Time `json:"admin"`
		// key: division_id value: 对应分区禁言解除时间
		Silent       map[int]*time.Time `json:"silent"`
		OffenseCount int                `json:"offense_count"`
	} `json:"permission" gorm:"-:all"`

	// get from jwt
	IsAdmin              bool      `json:"is_admin" gorm:"-:all"`
	JoinedTime           time.Time `json:"joined_time" gorm:"-:all"`
	Nickname             string    `json:"nickname" gorm:"-:all"`
	HasAnsweredQuestions bool      `json:"has_answered_questions" gorm:"-:all"`
}

type Users []*User

type UserConfig struct {
	// used when notify
	Notify []string `json:"notify"`

	// 对折叠内容的处理
	// fold 折叠, hide 隐藏, show 展示
	ShowFolded string `json:"show_folded"`
}

var defaultUserConfig = UserConfig{
	Notify:     []string{"mention", "favorite", "report"},
	ShowFolded: "hide",
}

var showFoldedOptions = []string{"hide", "fold", "show"}

func (user *User) GetID() int {
	return user.ID
}

func (user *User) AfterCreate(_ *gorm.DB) error {
	user.UserID = user.ID
	return nil
}

func (user *User) AfterFind(_ *gorm.DB) error {
	user.UserID = user.ID
	return nil
}

var (
	maxTime time.Time
	minTime time.Time
)

func init() {
	var err error
	maxTime, err = time.Parse(time.RFC3339, "9999-01-01T00:00:00+00:00")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	minTime = time.Unix(0, 0)
}

// GetCurrLoginUser get current login user
// In dev or test mode, return a default admin user
func GetCurrLoginUser(c *fiber.Ctx) (*User, error) {
	user := &User{
		BanDivision: make(map[int]*time.Time),
	}
	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		user.ID = 1
		user.IsAdmin = true
		user.HasAnsweredQuestions = true
		return user, nil
	}

	if c.Locals("user") != nil {
		return c.Locals("user").(*User), nil
	}

	// get id
	userID, err := common.GetUserID(c)
	if err != nil {
		return nil, err
	}

	// parse JWT
	err = common.ParseJWTToken(common.GetJWTToken(c), user)
	if err != nil {
		return nil, err
	}

	// load user from database in transaction
	err = user.LoadUserByID(userID)

	if user.IsAdmin {
		user.Permission.Admin = maxTime
	} else {
		user.Permission.Admin = minTime
	}
	user.Permission.Silent = user.BanDivision
	user.Permission.OffenseCount = user.OffenceCount

	if config.Config.UserAllShowHidden {
		user.Config.ShowFolded = "hide"
	}

	// save user in c.Locals
	c.Locals("user", user)

	return user, err
}

func (user *User) LoadUserByID(userID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Take(&user, userID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// insert user if not found
				user.ID = userID
				user.Config = defaultUserConfig
				err = tx.Create(&user).Error
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		err = CheckDefaultFavoriteGroup(tx, userID)
		if err != nil {
			return err
		}

		// check latest permission
		modified := false
		for divisionID := range user.BanDivision {
			endTime := user.BanDivision[divisionID]
			if endTime != nil && endTime.Before(time.Now()) {
				delete(user.BanDivision, divisionID)
				modified = true
			}
		}

		// check config
		if !slices.Contains(showFoldedOptions, user.Config.ShowFolded) {
			user.Config.ShowFolded = defaultUserConfig.ShowFolded
			modified = true
		}

		if user.Config.Notify == nil {
			user.Config.Notify = defaultUserConfig.Notify
			modified = true
		}

		if modified {
			err = tx.Select("BanDivision", "Config").Save(&user).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (user *User) BanDivisionMessage(divisionID int) string {
	if user.BanDivision[divisionID] == nil {
		return fmt.Sprintf("您在此板块已被禁言")
	} else {
		return fmt.Sprintf(
			"您在此板块已被禁言，解封时间：%s",
			user.BanDivision[divisionID].Format("2006-01-02 15:04:05"))
	}
}

func (user *User) BanReportMessage() string {
	if user.BanReport == nil {
		return fmt.Sprintf("您已被限制使用举报功能")
	} else {
		return fmt.Sprintf(
			"您已被限制使用举报功能，解封时间：%s",
			user.BanReport.Format("2006-01-02 15:04:05"))
	}
}
