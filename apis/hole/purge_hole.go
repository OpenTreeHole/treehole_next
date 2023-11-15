package hole

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"treehole_next/config"
	. "treehole_next/models"
)

func purgeHole() (err error) {
	const REASON = "purge_hole"
	const DELETE_CONTENT = "该内容已被删除"

	return DB.Transaction(func(tx *gorm.DB) (err error) {

		// load holeIDs, lock for update
		var holeIDs []int
		err = tx.Model(&Hole{}).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("no_purge = ?", false).
			Where("division_id IN ?", config.Config.HolePurgeDivisions).
			Where("updated_at < ?",
				time.Now().AddDate(0, 0, -config.Config.HolePurgeDays),
			).Pluck("id", &holeIDs).Error
		if err != nil {
			return err
		}

		if len(holeIDs) == 0 {
			return nil
		}

		/* delete all floors in hole of holeIOs */

		// get floors, lock for update
		var floors []Floor
		err = tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("hole_id IN ?", holeIDs).
			Find(&floors).Error
		if err != nil {
			return err
		}
		if len(floors) == 0 {
			return nil
		}

		// generate floorHistory
		var floorHistorySlice = make([]FloorHistory, 0, len(floors))
		for i := range floors {
			floorHistorySlice = append(floorHistorySlice, FloorHistory{
				Content: floors[i].Content,
				Reason:  REASON,
				FloorID: floors[i].ID,
				UserID:  1,
			})
		}
		err = tx.Create(&floorHistorySlice).Error
		if err != nil {
			return err
		}

		// delete floors
		var floorIDs = make([]int, 0, len(floors))
		for i := range floors {
			floorIDs = append(floorIDs, floors[i].ID)
		}

		err = tx.Model(&Floor{}).
			Where("id IN ?", floorIDs).
			Updates(map[string]any{
				"deleted": true,
				"content": DELETE_CONTENT,
			}).Error
		if err != nil {
			return err
		}

		/* delete all holes in holeIDs */
		err = tx.
			Where("id IN ?", holeIDs).
			Delete(&Hole{}).Error
		if err != nil {
			return err
		}

		// delete floor in search engine
		go BulkDelete(floorIDs)

		// log
		log.Info().
			Ints("hole_ids", holeIDs).
			Ints("floor_ids", floorIDs).
			Msg("purge hole")

		return nil
	})
}

func PurgeHole(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := purgeHole()
			if err != nil {
				log.Err(err).Msg("error purge hole")
			}
		case <-ctx.Done():
			return
		}
	}
}
