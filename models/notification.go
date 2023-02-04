package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
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
	Type        MessageType `json:"type" validate:"required"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Data        any         `json:"data" gorm:"serializer:json" `
	URL         string      `json:"url"`
	Recipients  []int       `json:"recipients" validate:"required"`
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
	messages = append(messages, newNotification)
	return messages
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

	// save to database
	body := Message{
		Type:        message.Type,
		Title:       utils.StripContent(message.Title, 32),       //varchar(32)
		Description: utils.StripContent(message.Description, 64), //varchar(64)
		Data:        message.Data,
		URL:         message.URL,
		Recipients:  message.Recipients,
	}
	err = DB.Create(&body).Error
	if err != nil {
		utils.Logger.Error("[notification] message save failed: " + err.Error())
		return Message{}, err
	}

	return body, nil
}

type Admin struct {
	Id           int      `json:"id"`
	IsAdmin      bool     `json:"is_admin"`
	JoinedTime   string   `json:"joined_time"`
	LastLogin    string   `json:"last_login"`
	Nickname     string   `json:"nickname"`
	OffenseCount int      `json:"offense_count"`
	Roles        []string `json:"roles"`
}

func readRespAdmin(body io.ReadCloser) []Admin {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			utils.Logger.Error("[get admin] Close error: " + err.Error())
		}
	}(body)

	data, err := io.ReadAll(body)
	if err != nil {
		utils.Logger.Error("[get admin] Read body failed: " + err.Error())
		return []Admin{}
	}
	var response []Admin
	err = json.Unmarshal(data, &response)
	if err != nil {
		utils.Logger.Error("[get admin] Unmarshal body failed: " + err.Error())
		return []Admin{}
	}
	return response
}

var adminList []int

func InitAdminList() {
	// skip when bench
	if config.Config.Mode == "bench" {
		return
	}

	// construct http request
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/users", config.Config.AuthUrl),
		nil,
	)
	query := req.URL.Query()
	query.Add("size", "0")
	query.Add("offset", "0")
	query.Add("role", "admin")
	req.URL.RawQuery = query.Encode()

	// get response
	resp, err := client.Do(req)

	// handle err
	if err != nil {
		utils.Logger.Error("[get admin] error sending auth server" + err.Error())
		return
	}

	response := readRespAdmin(resp.Body)

	if resp.StatusCode != 200 || len(response) == 0 {
		utils.Logger.Error("[get admin] auth server response failed" + fmt.Sprint(resp))
		return
	}

	// get ids
	for _, mention := range response {
		adminList = append(adminList, mention.Id)
	}

	// shuffle ids
	for i := range adminList {
		j := rand.Intn(i + 1)
		adminList[i], adminList[j] = adminList[j], adminList[i]
	}

	// panic(fmt.Sprint(adminList)) // Only for Test
}
