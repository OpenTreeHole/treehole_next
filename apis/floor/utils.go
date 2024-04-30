package floor

import "fmt"

func generateDeleteReason(reason string, isOwner bool) string {
	if reason == "" {
		if isOwner {
			return "该内容被作者删除"
		}
		reason = "违反社区规范"
	}
	return fmt.Sprintf("该内容因%s被删除", reason)
}
