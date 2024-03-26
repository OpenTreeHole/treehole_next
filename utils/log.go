package utils

import (
	"github.com/rs/zerolog/log"
	"treehole_next/utils/sensitive"
)

type Role string

const (
	RoleOwner    = "owner"
	RoleAdmin    = "admin"
	RoleOperator = "operator"
)

func MyLog(model string, action string, objectID, userID int, role Role, msg ...string) {
	message := ""
	for _, v := range msg {
		message += v
	}
	log.Info().
		Str("model", model).
		Int("user_id", userID).
		Int("object_id", objectID).
		Str("action", action).
		Str("role", string(role)).
		Msg(message)
}

func RequestLog(msg string, params sensitive.ParamsForCheck, ans bool) {
	log.Info().Str("TypeName", params.TypeName).
		Int("Id", params.Id).
		Bool("CheckAnswer", ans).
		Msg(msg)
}
