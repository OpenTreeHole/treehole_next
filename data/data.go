package data

import (
	_ "embed"
	"os"

	"github.com/goccy/go-json"
	"github.com/importcjj/sensitive"
	"github.com/rs/zerolog/log"
)

//go:embed names.json
var NamesFile []byte

//go:embed meta.json
var MetaFile []byte

var NamesMapping map[string]string

var SensitiveWordFilter *sensitive.Filter

func init() {
	err := initNamesMapping()
	if err != nil {
		log.Err(err).Msg("could not init names mapping")
	}

	err = initSensitiveWords()
	if err != nil {
		log.Err(err).Msg("could not init sensitive words")
	}
}

func initNamesMapping() error {
	NamesMappingData, err := os.ReadFile(`data/names_mapping.json`)
	if err != nil {
		return err
	}

	return json.Unmarshal(NamesMappingData, &NamesMapping)
}

func initSensitiveWords() error {
	SensitiveWordFilter = sensitive.New()
	err := SensitiveWordFilter.LoadWordDict("data/sensitive_words.txt")
	if err != nil {
		SensitiveWordFilter = nil
		return err
	}
	return nil
}
