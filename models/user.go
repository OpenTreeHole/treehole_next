package models

import (
	"encoding/base64"
	"errors"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"
	"strconv"
	"strings"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
)

type User struct {
	/// base info
	ID int `json:"id" gorm:"primaryKey"`

	Config UserConfig `json:"config" gorm:"serializer:json;not null;default:\"{}\""`

	BanDivision map[int]*time.Time `json:"-" gorm:"serializer:json;not null;default:\"{}\""`

	OffenceCount int `json:"-" gorm:"not null;default:0"`

	/// association fields, should add foreign key

	// favorite holes of the user
	UserFavoriteHoles Holes `json:"-" gorm:"many2many:user_favorite"`

	// holes owned by the user
	UserHoles Holes `json:"-"`

	// floors owned by the user
	UserFloors Floors `json:"-"`

	// reports made by the user; a user has many report
	UserReports Reports `json:"-"`

	// reports dealt by the user, admin only
	UserDealtReports Reports `json:"-" gorm:"foreignKey:DealtBy"`

	// floors liked by the user
	UserLikedFloors Floors `json:"-" gorm:"many2many:floor_like"`

	// floors disliked by the user
	UserDislikedFloors Floors `json:"-" gorm:"many2many:floor_dislike"`

	// floor history made by the user
	UserFloorHistory FloorHistorySlice `json:"-"`

	// user punishments on division
	UserPunishments Punishments `json:"-"`

	// punishments made by this user
	UserMakePunishments Punishments `json:"-" gorm:"foreignKey:MadeBy"`

	/// dynamically generated field

	Permission struct {
		// 管理员权限到期时间
		Admin time.Time `json:"admin"`
		// key: division_id value: 对应分区禁言解除时间
		Silence      map[int]*time.Time `json:"silence"`
		OffenseCount int                `json:"offense_count"`
	} `json:"permission" gorm:"-:all"`

	// load from table 'user_favorite'
	FavoriteData []int `json:"favorite_data" gorm:"-:all"`

	// get from jwt
	IsAdmin    bool      `json:"is_admin" gorm:"-:all"`
	JoinedTime time.Time `json:"joined_time" gorm:"-:all"`
	LastLogin  time.Time `json:"last_login" gorm:"-:all"`
	Nickname   string    `json:"nickname" gorm:"-:all"`
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
	ShowFolded: "fold",
}

func (user *User) GetID() int {
	return user.ID
}

// parseJWT extracts and parse token
func (user *User) parseJWT(token string) error {
	if len(token) < 7 {
		return errors.New("bearer token required")
	}

	payloads := strings.SplitN(token[7:], ".", 3) // extract "Bearer "
	if len(payloads) < 3 {
		return errors.New("jwt token required")
	}

	// jwt encoding ignores padding, so RawStdEncoding should be used instead of StdEncoding
	payloadBytes, err := base64.RawStdEncoding.DecodeString(payloads[1]) // the middle one is payload
	if err != nil {
		return err
	}

	err = json.Unmarshal(payloadBytes, user)
	if err != nil {
		return err
	}

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
		panic(err)
	}
	minTime = time.Unix(0, 0)
}

func GetUser(c *fiber.Ctx) (*User, error) {
	user := &User{
		BanDivision: make(map[int]*time.Time),
	}
	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		user.ID = 1
		user.IsAdmin = true
		return user, nil
	}

	// get id
	userID, err := GetUserID(c)
	if err != nil {
		return nil, err
	}

	// parse JWT first
	tokenString := c.Get("Authorization")
	if tokenString == "" { // token can be in either header or cookie
		tokenString = c.Cookies("access")
	}
	err = user.parseJWT(tokenString)
	if err != nil {
		return nil, utils.Unauthorized(err.Error())
	}

	// load user from database in transaction
	err = user.LoadUserByID(userID)

	if user.IsAdmin {
		user.Permission.Admin = maxTime
	} else {
		user.Permission.Admin = minTime
	}
	user.Permission.Silence = user.BanDivision
	user.Permission.OffenseCount = user.OffenceCount
	return user, err
}

func GetUserID(c *fiber.Ctx) (int, error) {
	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		return 1, nil
	}

	id, err := strconv.Atoi(c.Get("X-Consumer-Username"))
	if err != nil {
		return 0, utils.Unauthorized("Unauthorized")
	}

	return id, nil
}

func (user *User) LoadUserByID(userID int) error {
	return DB.Clauses(dbresolver.Write).Transaction(func(tx *gorm.DB) error {
		err := tx.Preload("UserPunishments").Clauses(clause.Locking{Strength: "UPDATE"}).Take(&user, userID).Error
		if err == gorm.ErrRecordNotFound {
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

		// check permission
		modified := false
		for divisionID := range user.BanDivision {
			// get the latest punishments in divisionID
			var latestPunishment *Punishment
			for _, punishment := range user.UserPunishments {
				if punishment.DivisionID == divisionID {
					latestPunishment = punishment
				}
			}

			if latestPunishment == nil || latestPunishment.EndTime.Before(time.Now()) {
				delete(user.BanDivision, divisionID)
				modified = true
			}
		}

		if modified {
			err = tx.Select("BanDivision").Save(&user).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (user *User) Preprocess(_ *fiber.Ctx) (err error) {
	user.FavoriteData, err = UserGetFavoriteData(DB, user.ID)
	return err
}
