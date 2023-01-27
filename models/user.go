package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
	"treehole_next/utils/perm"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	/// base info
	ID int `json:"id" gorm:"primaryKey"`

	Config struct {
		// used when notify
		Notify []string `json:"notify"`

		// 对折叠内容的处理
		// fold 折叠, hide 隐藏, show 展示
		ShowFolded string `json:"show_folded"`
	} `json:"config" gorm:"serializer:json;not null;default:\"{}\""`

	/// association fields, should add foreign key

	// favorite holes of the user
	UserFavoriteHoles []*Hole `json:"-" gorm:"many2many:user_favorite"`

	// holes owned by the user
	UserHoles []Hole `json:"-"`

	// floors owned by the user
	UserFloors []Floor `json:"-"`

	// reports made by the user; a user has many report
	UserReports []Report `json:"-"`

	// reports dealt by the user, admin only
	UserDealtReports []Report `json:"-" gorm:"foreignKey:DealtBy"`

	// floors liked by the user
	UserLikedFloors []*Floor `json:"-" gorm:"many2many:floor_like"`

	// floors disliked by the user
	UserDislikedFloors []*Floor `json:"-" gorm:"many2many:floor_dislike"`

	// floor history made by the user
	UserFloorHistory []FloorHistory `json:"-"`

	// user punishments on division
	UserPunishments []Punishment `json:"-"`

	// punishments made by this user
	UserMakePunishments []Punishment `json:"-" gorm:"foreignKey:MadeBy"`

	/// dynamically generated field

	Permission struct {
		// 管理员权限到期时间
		Admin time.Time `json:"admin"`
		// key: division_id value: 对应分区禁言解除时间
		Silence      map[int]time.Time `json:"silence"`
		OffenseCount int               `json:"offense_count"`
	} `json:"permission" gorm:"-:all"`

	BanDivision map[int]bool `json:"-" gorm:"-:all"`

	// load from table 'user_favorite'
	FavoriteData []int `json:"favorite_data" gorm:"-:all"`

	// get from jwt
	IsAdmin    bool      `json:"is_admin" gorm:"-:all"`
	JoinedTime time.Time `json:"joined_time" gorm:"-:all"`
	LastLogin  time.Time `json:"last_login" gorm:"-:all"`
	Nickname   string    `json:"nickname" gorm:"-:all"`

	// deprecated
	Role perm.Role `json:"-" gorm:"-:all"`
}

func (user User) GetID() int {
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

func GetUser(c *fiber.Ctx) (*User, error) {
	user := &User{
		BanDivision: make(map[int]bool),
		Role:        0,
	}
	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		user.ID = 1
		user.Role = perm.Admin + perm.Operator
		return user, nil
	}

	// get id
	id, err := GetUserID(c)
	if err != nil {
		return nil, err
	}
	user.ID = id

	// parse JWT
	tokenString := c.Get("Authorization")
	if tokenString == "" { // token can be in either header or cookie
		tokenString = c.Cookies("access")
	}
	err = user.parseJWT(tokenString)
	if err != nil {
		return nil, utils.Unauthorized(err.Error())
	}

	err = user.parsePermission()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (user *User) parsePermission() error {
	if user.IsAdmin {
		user.Role |= perm.Admin
	}
	return nil
}

func GetUserFromAuth(c *fiber.Ctx) (*User, error) {
	user := &User{
		BanDivision: make(map[int]bool),
		Role:        0,
	}

	if config.Config.Mode == "dev" || config.Config.Mode == "test" {
		user.ID = 1
		user.Role = perm.Admin + perm.Operator
		return user, nil
	}

	userID, err := GetUserID(c)
	if err != nil {
		return nil, err
	}

	// make request
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/users/%d", config.Config.AuthUrl, userID),
		bytes.NewBuffer(make([]byte, 0, 10)),
	)
	if err != nil {
		utils.Logger.Error("request make error", zap.Error(err))
		return nil, err
	}

	// add headers
	req.Header.Add("X-Consumer-Username", strconv.Itoa(userID))
	rsp, err := client.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Logger.Error("close rsp body error")
		}
	}(rsp.Body)

	if err != nil {
		utils.Logger.Error(
			"auth get user request error",
			zap.Int("user id", userID),
		)
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, errors.New("auth get user error, rsp error")
	}

	userInfo, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(userInfo, user)
	if err != nil {
		return nil, err
	}

	err = user.parsePermission()
	if err != nil {
		return nil, err
	}

	return user, nil
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

func (user *User) GetPermission() perm.Role {
	return user.Role
}
