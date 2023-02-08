package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"treehole_next/config"
	"treehole_next/utils"

	"github.com/goccy/go-json"
)

const (
	timeout = time.Second * 10
)

var client = http.Client{Timeout: timeout}

type Notifications []Notification

type Notification struct {
	// Should be same as CrateModel in notification project
	Title       string      `json:"message"`
	Description string      `json:"description"`
	Data        any         `json:"data"`
	Type        MessageType `json:"code"`
	URL         string      `json:"url"`
	Recipients  []int       `json:"recipients"`
}

func readRespNotification(body io.ReadCloser) Notification {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			utils.Logger.Error("[notification] Close error: " + err.Error())
		}
	}(body)

	data, err := io.ReadAll(body)
	if err != nil {
		utils.Logger.Error("[notification] Read body failed: " + err.Error())
		return Notification{}
	}
	var response Notification
	err = json.Unmarshal(data, &response)
	if err != nil {
		utils.Logger.Error("[notification] Unmarshal body failed: " + err.Error())
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

func (message Notification) Send() (Message, error) {
	// only for test
	// message["recipients"] = []int{1}

	var err error

	// save to database first
	body := Message{
		Type:        message.Type,
		Title:       utils.StripContent(message.Title, 32),       //varchar(32)
		Description: utils.StripContent(message.Description, 64), //varchar(64)
		Data:        message.Data,
		URL:         message.URL,
		Recipients:  message.Recipients,
	}
	err = DB.Omit(clause.Associations).Create(&body).Error
	if err != nil {
		log.Println("[notification] message save failed: " + err.Error())
		return Message{}, err
	}
	if config.Config.NotificationUrl == "" {
		return Message{}, nil
	}

	// construct form
	form, err := json.Marshal(message)
	if err != nil {
		utils.Logger.Error("[notification] error encoding notification: " + err.Error())
		return Message{}, err
	}

	// construct http request
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/messages", config.Config.NotificationUrl),
		bytes.NewBuffer(form),
	)
	if err != nil {
		utils.Logger.Error("[notification] error making request: " + err.Error())
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
		utils.Logger.Error("[notification] error sending notification: " + err.Error())
		return Message{}, err
	}

	response := readRespNotification(resp.Body)
	if resp.StatusCode != 201 {
		utils.Logger.Error("[notification] notification response failed: " + fmt.Sprint(response))
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

	// http request
	res, err := http.Get(config.Config.AuthUrl + "/users/admin")

	// handle err
	if err != nil {
		utils.Logger.Error("[get admin] error sending auth server" + err.Error())
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		utils.Logger.Error("[get admin] auth server response failed" + res.Status)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		utils.Logger.Error("[get admin] auth server response failed" + err.Error())
		return
	}

	adminList.Lock()
	defer adminList.Unlock()

	err = json.Unmarshal(data, &adminList.data)
	if err != nil {
		utils.Logger.Error("[get admin] auth server response failed" + err.Error())
		return
	}

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
