package models

import (
	"fmt"
	"treehole_next/utils"
)

const contentMaxSize = 50

func getContent(data Map) string {
	content, ok := data["content"].(string)
	if !ok {
		return ""
	}
	return content
}

func getFloorContent(data Map) string {
	floor, ok := data["floor"].(Map)
	if !ok {
		return ""
	}
	return getContent(floor)
}

func stripContent(content string) string {
	return content[:utils.Min(len(content), contentMaxSize)]
}

func generateTitle(messageType MessageType) string {
	switch messageType {
	case MessageTypeFavorite:
		return "你收藏的树洞有新回复"
	case MessageTypeReply:
		return "你的帖子被回复了"
	case MessageTypeMention:
		return "你的帖子被引用了"
	case MessageTypeModify:
		return "你的帖子被修改了"
	case MessageTypePermission:
		return "你的权限被更改了"
	case MessageTypeReport:
		return "有帖子被举报了"
	case MessageTypeReportDealt:
		return "你的举报被处理了"
	}
	return "通知"
}

func generateDescription(messageType MessageType, data Map) string {
	switch messageType {
	// data is floor
	case MessageTypeFavorite:
		return getContent(data)
	case MessageTypeReply:
		return getContent(data)
	case MessageTypeMention:
		return getContent(data)
	case MessageTypeModify:
		return getContent(data)
	// data is permission
	case MessageTypePermission:
		return fmt.Sprintf(
			"权限：%s，理由：%s，截止时间：%s",
			data["name"],
			data["reason"],
			data["end_time"],
		)
	// data is report
	case MessageTypeReport:
		return fmt.Sprintf(
			"内容：%s，理由：%s",
			stripContent(getFloorContent(data)),
			data["reason"],
		)
	case MessageTypeReportDealt:
		return fmt.Sprintf(
			"内容：%s，结果：%s",
			stripContent(getFloorContent(data)),
			data["result"],
		)
	}
	return ""
}
