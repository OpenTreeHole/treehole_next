package utils

import (
	"github.com/rs/zerolog/log"
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
