package data

import (
	_ "embed"
)

//go:embed names.json
var NamesFile []byte

//go:embed meta.json
var MetaFile []byte
