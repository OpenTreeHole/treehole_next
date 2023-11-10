package utils

import (
	"treehole_next/config"
	"treehole_next/data"
)

func IsSensitive(content string) bool {
	if !config.Config.OpenSensitiveCheck {
		return false
	}

	if data.SensitiveWordFilter != nil {
		in, _ := data.SensitiveWordFilter.FindIn(content)
		return in
	}
	return true
}
