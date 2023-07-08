package utils

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sort"
	"time"
	"treehole_next/config"
	"treehole_next/data"

	"golang.org/x/exp/slices"
)

var names []string
var length int

func init() {
	err := json.Unmarshal(data.NamesFile, &names)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	sort.Strings(names)
	length = len(names)
}

func inArray(target string, array []string) bool {
	_, in := slices.BinarySearch(array, target)
	return in
}

func timeStampBase64() string {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(time.Now().Unix()))
	return base64.StdEncoding.EncodeToString(bytes)
}

func NewRandName() string {
	return names[rand.Intn(length)]
}

func GenerateName(compareList []string) string {
	if len(compareList) < length>>3 {
		for {
			name := NewRandName()
			if !inArray(name, compareList) {
				return name
			}
		}
	} else if len(compareList) < length {
		var j, k int
		list := make([]string, length)
		for i := 0; i < length; i++ {
			if j < len(compareList) && names[i] == compareList[j] {
				j++
			} else {
				list[k] = names[i]
				k++
			}
		}
		return list[rand.Intn(k)]
	} else {
		for {
			name := names[rand.Intn(length)] + "_" + timeStampBase64()
			if !inArray(name, compareList) {
				return name
			}
		}
	}
}

func GetFuzzName(name string) string {
	if !config.Config.OpenFuzzName {
		return name
	}
	if fuzzName, ok := data.NamesMapping[name]; ok {
		return fuzzName
	} else {
		return name
	}
}
