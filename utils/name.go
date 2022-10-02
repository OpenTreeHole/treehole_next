package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"sort"
	"time"
	"treehole_next/data"

	"golang.org/x/exp/slices"
)

var names []string
var length int

func init() {
	err := json.Unmarshal(data.NamesFile, &names)
	if err != nil {
		panic(err)
	}
	sort.Strings(names)
	length = len(names)

	rand.Seed(time.Now().UnixNano())
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

func GenerateName(compareList []string) string {
	if len(compareList) < length>>3 {
		for {
			name := names[rand.Intn(length)]
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
