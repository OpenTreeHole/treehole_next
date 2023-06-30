package division

import (
	"github.com/rs/zerolog/log"
	. "treehole_next/models"
)

func refreshCache() {
	var divisions Divisions
	DB.Find(&divisions)
	err := divisions.Preprocess(nil)
	if err != nil {
		log.Err(err).Msg("error refreshing cache")
	}
}
