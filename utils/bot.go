package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"treehole_next/config"
)

type BotMessageType string

const (
	MessageTypeGroup   BotMessageType = "group"
	MessageTypePrivate BotMessageType = "private"
)

type BotMessage struct {
	MessageType BotMessageType `json:"message_type"`
	GroupID     *int64         `json:"group_id"`
	UserID      *int64         `json:"user_id"`
	Message     string         `json:"message"`
	AutoEscape  bool           `json:"auto_escape default:false"`
}

func NotifyQQ(botMessage *BotMessage) {
	if botMessage == nil {
		return
	}
	if botMessage.MessageType == MessageTypeGroup && botMessage.GroupID == nil {
		return
	}
	if botMessage.MessageType == MessageTypePrivate && botMessage.UserID == nil {
		return
	}
	if config.Config.QQBotUrl == nil {
		return
	}
	// "[CQ:face,id=199]test[CQ:image,file=https://ts1.cn.mm.bing.net/th?id=OIP-C.K5AFHsGlWeLUzKjXGXxdQgHaFj&w=224&h=150&c=8&rs=1&qlt=90&o=6&dpr=1.5&pid=3.1&rm=2]",
	url := *config.Config.QQBotUrl + "/send_msg"

	jsonData, err := json.Marshal(botMessage)
	if err != nil {
		RequestLog("Error marshaling JSON", "NotifyQQ", 0, false)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		RequestLog("Error creating request", "NotifyQQ", 0, false)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		response, err := io.ReadAll(resp.Body)
		if err != nil {
			RequestLog("Error Unmarshaling response", "NotifyQQ", 0, false)
		}
		RequestLog(fmt.Sprintf("Error sending request %s", string(response)), "NotifyQQ", 0, false)
		request, _ := json.Marshal(botMessage)
		RequestLog(fmt.Sprintf("Request: %s", string(request)), "NotifyQQ", 0, false)
	}
}
