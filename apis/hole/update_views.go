package hole

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
	. "treehole_next/models"
)

var holeViewsChan = make(chan int, 1000)
var holeViews = map[int]int{}

func updateHoleViews() {
	/*
		UPDATE table
		SET field = CASE id
			WHEN 1 THEN 'value'
			WHEN 2 THEN 'value'
			WHEN 3 THEN 'value'
		END
		WHERE id IN (1,2,3)
	*/
	length := len(holeViews)
	if length == 0 {
		return
	}
	keys := make([]string, 0, length)

	var builder strings.Builder
	builder.WriteString("UPDATE hole SET view = CASE id ")

	for holeID, views := range holeViews {
		builder.WriteString(fmt.Sprintf("WHEN %d THEN view + %d ", holeID, views))
		keys = append(keys, strconv.Itoa(holeID))
		delete(holeViews, holeID)
	}
	builder.WriteString("END WHERE id IN (")
	builder.WriteString(strings.Join(keys, ","))
	builder.WriteString(")")

	result := DB.Exec(builder.String())
	if result.Error != nil {
		log.Err(result.Error).Msg("update hole views failed")
	} else {
		log.Info().Strs("updated", keys).Msg("update hole views success")
	}
}

func UpdateHoleViews(ctx context.Context) {

	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			updateHoleViews()
		case holeID := <-holeViewsChan:
			holeViews[holeID]++
		case <-ctx.Done():
			updateHoleViews()
			log.Info().Msg("task UpdateHoleViews stopped...")
			return
		}
	}
}
