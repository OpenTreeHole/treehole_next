package division

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	. "treehole_next/models"
)

func refreshCache(c *fiber.Ctx) error {

	var divisions Divisions
	err := DB.Find(&divisions).Error
	if err != nil {
		return err
	}

	err = divisions.Preprocess(c)
	if err != nil {
		log.Err(err).Msg("error refreshing cache")
		return err
	}

	return nil
}
