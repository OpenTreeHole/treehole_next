package data

import (
	_ "embed"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"os"
)

//go:embed names.json
var NamesFile []byte

//go:embed meta.json
var MetaFile []byte

var NamesMapping map[string]string

func init() {
	NamesMappingData, err := os.ReadFile(`data/names_mapping.json`)
	if err != nil {
		log.Err(err).Msg("could not load names_mapping.json")
		return
	}

	err = json.Unmarshal(NamesMappingData, &NamesMapping)
	if err != nil {
		log.Err(err).Msg("could not unmarshal names_mapping.json")
		return
	}
}
