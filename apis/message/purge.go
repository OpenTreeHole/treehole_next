package message

import (
	"time"

	"github.com/rs/zerolog/log"

	"treehole_next/config"
	. "treehole_next/models"
)

func purgeMessage() error {
	return DB.Exec(
		"DELETE FROM message WHERE created_at < ?",
		time.Now().Add(-time.Hour*24*time.Duration(config.Config.MessagePurgeDays)),
	).Error
}

func PurgeMessage() {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	for range ticker.C {
		err := purgeMessage()
		if err != nil {
			log.Err(err).Msg("error purge message")
		}
	}
}
