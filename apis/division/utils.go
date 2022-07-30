package division

import (
	"go.uber.org/zap"
	. "treehole_next/models"
	"treehole_next/utils"
)

func refreshCache() {
	var divisions Divisions
	DB.Find(&divisions)
	err := divisions.Preprocess(nil)
	if err != nil {
		utils.Logger.Error("error refreshing cache", zap.String("error", err.Error()))
	}
}
