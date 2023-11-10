package data

import (
	_ "embed"
	"os"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

//go:embed names.json
var NamesFile []byte

//go:embed meta.json
var MetaFile []byte

var NamesMapping map[string]string

var SensitiveWords []string

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
	sensitiveWordsData, err := os.ReadFile(`data/sensitive_words.json`)
	if err != nil {
		return err
	}

	return json.Unmarshal(sensitiveWordsData, &SensitiveWords)
}
