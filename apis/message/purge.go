package message

import (
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"treehole_next/config"
	. "treehole_next/models"
)

func purgeMessage() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// delete outdated messages
		result := tx.Exec(
			"DELETE FROM message WHERE created_at < ?",
			time.Now().Add(-time.Hour*24*time.Duration(config.Config.MessagePurgeDays)),
		)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
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
