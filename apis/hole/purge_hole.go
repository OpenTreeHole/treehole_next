package hole

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"treehole_next/config"
	. "treehole_next/models"
)

func purgeHole() error {
	return DB.
		Where("no_purge = ?", false).
		Where("division_id IN ?", config.Config.HolePurgeDivisions).
		Where("updated_at < ?",
			time.Now().AddDate(0, 0, -config.Config.HolePurgeDays),
		).Delete(&Hole{}).Error
}

func PurgeHole(ctx context.Context) {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := purgeHole()
			if err != nil {
				log.Err(err).Msg("error purge message")
			}
		case <-ctx.Done():
			return
		}
	}
}
