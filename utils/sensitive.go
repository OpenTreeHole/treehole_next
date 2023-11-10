package utils

import (
	"treehole_next/config"
	"treehole_next/data"
)

func IsSensitive(content string, weak bool) bool {
	if !config.Config.OpenSensitiveCheck {
		return false
	}

	if weak {
		if data.WeakSensitiveWordFilter != nil {
			in, _ := data.WeakSensitiveWordFilter.FindIn(content)
			return in
		}
	} else {
		if data.SensitiveWordFilter != nil {
			in, _ := data.SensitiveWordFilter.FindIn(content)
			return in
		}
	}

	return true
}
