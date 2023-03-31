package data

import (
	_ "embed"
	"github.com/goccy/go-json"
	"log"
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
		log.Println(err)
		return
	}

	err = json.Unmarshal(NamesMappingData, &NamesMapping)
	if err != nil {
		log.Println(err)
		return
	}
}
