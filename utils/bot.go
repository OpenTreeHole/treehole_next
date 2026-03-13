package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"treehole_next/config"
)

type NotificationTarget string

const (
	NotificationTargetQQUser         NotificationTarget = "qq_user"
	NotificationTargetQQPhysicsGroup NotificationTarget = "qq_physics_group"
	NotificationTargetQQCodingGroup  NotificationTarget = "qq_coding_group"
	NotificationTargetFeishuAdmin    NotificationTarget = "feishu_admin"
	NotificationTargetFeishuDivision NotificationTarget = "feishu_division"
)

type Notifier interface {
	Notify(target NotificationTarget, message string)
}

type BotMessageType string

const (
	MessageTypeGroup   BotMessageType = "group"
	MessageTypePrivate BotMessageType = "private"
)

type qqBotMessage struct {
	MessageType BotMessageType `json:"message_type"`
	GroupID     *int64         `json:"group_id"`
	UserID      *int64         `json:"user_id"`
	Message     string         `json:"message"`
	AutoEscape  bool           `json:"auto_escape default:false"`
}

type feishuMessage struct {
	MsgType string `json:"msg_type"`
	Content string `json:"message"`
}

type botNotifier struct{}

var defaultNotifier Notifier = botNotifier{}
var notificationHTTPClient = &http.Client{Timeout: 5 * time.Minute}

func Notify(target NotificationTarget, message string) {
	defaultNotifier.Notify(target, message)
}

func (botNotifier) Notify(target NotificationTarget, message string) {
	if message == "" {
		return
	}

	go func() {
		switch target {
		case NotificationTargetQQUser:
			notifyQQ(&qqBotMessage{
				MessageType: MessageTypePrivate,
				UserID:      config.Config.QQBotUserID,
				Message:     message,
			})
		case NotificationTargetQQPhysicsGroup:
			notifyQQ(&qqBotMessage{
				MessageType: MessageTypeGroup,
				GroupID:     config.Config.QQBotPhysicsGroupID,
				Message:     message,
			})
		case NotificationTargetQQCodingGroup:
			notifyQQ(&qqBotMessage{
				MessageType: MessageTypeGroup,
				GroupID:     config.Config.QQBotCodingGroupID,
				Message:     message,
			})
		case NotificationTargetFeishuAdmin:
			notifyFeishu(config.Config.FeishuAdminNotifierUrl, &feishuMessage{
				MsgType: "text",
				Content: message,
			})
		case NotificationTargetFeishuDivision:
			notifyFeishu(config.Config.FeishuDivisionNotifierUrl, &feishuMessage{
				MsgType: "text",
				Content: message,
			})
		}
	}()
}

func notifyFeishu(url *string, feishuMessage *feishuMessage) {
	if feishuMessage == nil || feishuMessage.MsgType == "" {
		return
	}
	if url == nil || *url == "" {
		return
	}

	jsonData, err := json.Marshal(feishuMessage)
	if err != nil {
		RequestLog("Error marshaling JSON", "NotifyFeishu", 0, false)
		return
	}

	RequestLog(fmt.Sprintf("Request: %s", string(jsonData)), "NotifyFeishu", 0, false)

	resp, err := notificationHTTPClient.Post(*url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		RequestLog("Error creating request", "NotifyFeishu", 0, false)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		response, err := io.ReadAll(resp.Body)
		if err != nil {
			RequestLog("Error Unmarshaling response", "NotifyFeishu", 0, false)
		}
		RequestLog(fmt.Sprintf("Error sending request %s", string(response)), "NotifyFeishu", 0, false)
	}
}

func notifyQQ(botMessage *qqBotMessage) {
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

	RequestLog(fmt.Sprintf("Request: %s", string(jsonData)), "NotifyQQ", 0, false)

	resp, err := notificationHTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
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
	}
}
