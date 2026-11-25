package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"treehole_next/config"
	"treehole_next/utils"

	"golang.org/x/exp/slices"
	"gorm.io/gorm/clause"

	"github.com/goccy/go-json"
)

const (
	timeout = time.Second * 10
)

var client = http.Client{Timeout: timeout}

type Notifications []Notification

type Notification struct {
	// Should be same as CrateModel in notification project
	Title          string      `json:"message"`
	Description    string      `json:"description"`
	Data           any         `json:"data"`
	Type           MessageType `json:"code"`
	URL            string      `json:"url"`
	Recipients     []int       `json:"recipients"`
	RelatedFloorID *int        `json:"related_floor_id,omitempty"`
	RelatedHoleID  *int        `json:"related_hole_id,omitempty"`
}

func readRespNotification(body io.ReadCloser) Notification {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Err(err).Str("model", "Notification").Msg("error close body")
		}
	}(body)

	data, err := io.ReadAll(body)
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("error read body")
		return Notification{}
	}
	var response Notification
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("error unmarshal body")
		return Notification{}
	}
	return response
}

func (messages Notifications) Merge(newNotification Notification) Notifications {
	if len(newNotification.Recipients) == 0 {
		return messages
	}

	newMerge := newNotification.Recipients
	for _, message := range messages {
		old := message.Recipients
		for _, r1 := range old {
			for id, r2 := range newMerge {
				if r1 == r2 {
					newMerge = append(newMerge[:id], newMerge[id+1:]...)
					break
				}
			}
		}
		if len(newMerge) == 0 {
			return messages
		}
	}

	newNotification.Recipients = newMerge
	return append(messages, newNotification)
}

func (messages Notifications) Send() error {
	if messages == nil {
		return nil
	}

	for _, message := range messages {
		_, err := message.Send()
		if err != nil {
			return err
		}
	}
	return nil
}

// check user.config.Notify contain message.Type
func (message *Notification) checkConfig() {
	// generate new recipients
	var newRecipient []int

	// find users
	var users []User
	result := DB.Find(&users, "id in ?", message.Recipients)
	if result.Error != nil {
		message.Recipients = newRecipient
		return
	}

	// filter recipients
	for _, user := range users {
		if slices.Contains(defaultUserConfig.Notify, string(message.Type)) && !slices.Contains(user.Config.Notify, string(message.Type)) {
			continue
		}
		newRecipient = append(newRecipient, user.ID)
	}
	message.Recipients = newRecipient
}

func (message Notification) Send() (Message, error) {
	// only for test
	// message["recipients"] = []int{1}

	var err error

	message.checkConfig()
	if len(message.Recipients) == 0 {
		return Message{}, nil
	}
	var restricted []int
	if err = DB.Model(&UserContentFilter{}).Where("user_id IN ?", message.Recipients).Pluck("user_id", &restricted).Error; err != nil {
		return Message{}, err
	}
	if len(restricted) == 0 {
		return message.deliver()
	}
	ordinary := make([]int, 0, len(message.Recipients))
	for _, id := range message.Recipients {
		if !slices.Contains(restricted, id) {
			ordinary = append(ordinary, id)
		}
	}
	filtered := message
	filtered.Recipients = nil
	filtered.Data = nil
	for _, id := range restricted {
		visible, err := message.visibleTo(id)
		if err != nil {
			return Message{}, err
		}
		if visible {
			filtered.Recipients = append(filtered.Recipients, id)
		}
	}
	message.Recipients = ordinary
	first, err := message.deliver()
	if err != nil {
		return first, err
	}
	second, err := filtered.deliver()
	if first.ID != 0 {
		return first, err
	}
	return second, err
}

func (message Notification) visibleTo(uid int) (bool, error) {
	if message.RelatedFloorID == nil && message.RelatedHoleID == nil {
		return message.Type == MessageTypeMail, nil
	}
	if message.RelatedFloorID != nil {
		var count int64
		err := VisibleFloorIDs(DB, uid).Where("floor.id = ?", *message.RelatedFloorID).Count(&count).Error
		if err != nil || count == 0 {
			return false, err
		}
	}
	if message.RelatedHoleID != nil {
		var count int64
		err := VisiblePublicHoleIDs(DB, uid).Where("hole.id = ?", *message.RelatedHoleID).Count(&count).Error
		if err != nil || count == 0 {
			return false, err
		}
	}
	return true, nil
}

func (message Notification) deliver() (Message, error) {
	var err error
	// return if no recipient
	if len(message.Recipients) == 0 {
		return Message{}, nil
	}

	// save to database first
	body := Message{
		Type:           message.Type,
		Title:          message.Title,
		Description:    message.Description,
		Data:           message.Data,
		URL:            message.URL,
		Recipients:     message.Recipients,
		RelatedFloorID: message.RelatedFloorID,
		RelatedHoleID:  message.RelatedHoleID,
	}
	err = DB.Omit(clause.Associations).Create(&body).Error
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("message save failed: " + err.Error())
		return Message{}, err
	}
	if config.Config.NotificationUrl == "" {
		return Message{}, nil
	}
	message.Title = utils.StripContent(message.Title, 32)                                           //varchar(32)
	message.Description = utils.StripContent(cleanNotificationDescription(message.Description), 64) //varchar(64)
	body.Title = message.Title
	body.Description = message.Description

	// construct form
	form, err := json.Marshal(message)
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("error encoding notification")
		return Message{}, err
	}

	// construct http request
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/messages", config.Config.NotificationUrl),
		bytes.NewBuffer(form),
	)
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("error making request")
		return Message{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	// bench and simulation
	if config.Config.Mode == "bench" {
		time.Sleep(time.Millisecond)
		return Message{}, nil
	}

	// get response
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Str("model", "Notification").Msg("error sending notification")
		return Message{}, err
	}

	response := readRespNotification(resp.Body)
	if resp.StatusCode != 201 {
		log.Error().Str("model", "Notification").Any("response", response).Msg("notification response failed")
		return Message{}, errors.New(fmt.Sprint(response))
	}

	return body, nil
}

var adminList struct {
	sync.RWMutex
	data []int
}

func InitAdminList() {
	// skip when bench
	if config.Config.Mode == "bench" || config.Config.AuthUrl == "" {
		return
	}

	// // http request
	// res, err := http.Get(config.Config.AuthUrl + "/users/admin")

	// // handle err
	// if err != nil {
	// 	log.Err(err).Str("model", "get admin").Msg("error sending auth server")
	// 	return
	// }

	// defer func() {
	// 	_ = res.Body.Close()
	// }()

	// if res.StatusCode != 200 {
	// 	log.Error().Str("model", "get admin").Msg("auth server response failed" + res.Status)
	// 	return
	// }

	// data, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Err(err).Str("model", "get admin").Msg("error reading auth server response")
	// 	return
	// }

	adminList.Lock()
	defer adminList.Unlock()

	// err = json.Unmarshal(data, &adminList.data)
	// if err != nil {
	// 	log.Err(err).Str("model", "get admin").Msg("error unmarshal auth server response")
	// 	return
	// }
	adminList.data = config.Config.NotifiableAdminIds

	// shuffle ids
	for i := range adminList.data {
		j := rand.Intn(i + 1)
		adminList.data[i], adminList.data[j] = adminList.data[j], adminList.data[i]
	}
}

func UpdateAdminList(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			InitAdminList()
		}
	}
}

var (
	reMention = regexp.MustCompile(`#{1,2}\d+`)
	reFormula = regexp.MustCompile(`(?s)\${1,2}.*?\${1,2}`)
	reSticker = regexp.MustCompile(`!\[\]\(dx_\S+?\)`)
	reImage   = regexp.MustCompile(`!\[.*?\]\(.*?\)`)
)

func cleanNotificationDescription(content string) string {
	newContent := reMention.ReplaceAllString(content, "")
	newContent = reFormula.ReplaceAllString(newContent, "[公式]")
	newContent = reSticker.ReplaceAllString(newContent, "[表情]")
	newContent = reImage.ReplaceAllString(newContent, "[图片]")
	newContent = strings.ReplaceAll(newContent, "\n", "")
	if newContent == "" {
		return content
	}
	return newContent
}
