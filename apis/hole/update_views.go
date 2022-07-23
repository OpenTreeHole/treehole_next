package hole

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	. "treehole_next/models"
	"treehole_next/utils"
)

var holeViewsChan = make(chan int, 100)

func receiveViewsUpdate() {
	for {
		select {
		case holeID := <-holeViewsChan:
			holeViews[holeID]++
		}
	}
}

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
		utils.Logger.Error(result.Error.Error())
	} else {
		utils.Logger.Info(
			"update hole views success",
			zap.Strings("updated", keys),
		)
	}
}

func UpdateHoleViews() {
	go receiveViewsUpdate()

	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for range ticker.C {
		updateHoleViews()
	}
}
