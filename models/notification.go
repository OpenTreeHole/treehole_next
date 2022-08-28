package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"treehole_next/config"
	"treehole_next/utils"
)

const (
	timeout = time.Second * 10
)

var client = http.Client{Timeout: timeout}

type JSON map[string]any

func (t JSON) Value() ([]byte, error) {
	return json.Marshal(t)
}

func (t *JSON) Scan(input any) error {
	return json.Unmarshal(input.([]byte), t)
}

/*
Message Define

	{
		"data": { // json of models
		"additionalProp1": "string",
		"additionalProp2": "string",
		"additionalProp3": "string"
		},
		"description": "string", // LEAVE BLANK, will generate by micro service
		"recipients": [ // UserId
		0
		],
		"title": "string", // LEAVE BLANK, will generate by micro service
		"type": "favorite", // define by MessageType
		"url": "string" // relative api route
	}
*/
type Message map[string]any

type MessageType string

const (
	MessageTypeFavorite    MessageType = "favorite"
	MessageTypeReply       MessageType = "reply"
	MessageTypeMention     MessageType = "mention"
	MessageTypeModify      MessageType = "modify" // including fold and delete
	MessageTypePermission  MessageType = "permission"
	MessageTypeReport      MessageType = "report"
	MessageTypeReportDealt MessageType = "report_dealt"
)

func readBody(body io.ReadCloser) Map {
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			utils.Logger.Error("[notification] Close error: " + err.Error())
		}
	}(body)

	data, err := ioutil.ReadAll(body)
	if err != nil {
		utils.Logger.Error("[notification] Read body failed: " + err.Error())
		return Map{}
	}
	var response Map
	err = json.Unmarshal(data, &response)
	if err != nil {
		utils.Logger.Error("[notification] Unmarshal body failed: " + err.Error())
		return Map{}
	}
	return response
}

func (message Message) Send() error {
	// only for testing
	// message["recipients"] = []int{1}

	// construct form
	form, err := json.Marshal(message)
	if err != nil {
		utils.Logger.Error("[notification] error encoding notification" + err.Error())
		return err
	}

	// construct http request
	req, _ := http.NewRequest(
		"POST",
		config.Config.MicroUrl,
		bytes.NewBuffer(form),
	)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	// get response and handle err
	response := readBody(resp.Body)
	if err != nil {
		utils.Logger.Error("[notification] error sending notification" + err.Error())
		return err
	} else if resp.StatusCode != 201 {
		utils.Logger.Error("[notification] microservice response failed")
		return errors.New(fmt.Sprint(response))
	}

	return nil
}
