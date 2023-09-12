package benchmarks

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	. "treehole_next/models"
	_ "treehole_next/tests"
)

func BenchmarkListHoles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// prepare
		b.StopTimer()
		route := "/api/divisions/" + strconv.Itoa(rand.Intn(DIVISION_MAX)+1) + "/holes/"
		b.StartTimer()

		benchmarkCommon(b, "get", route, REQUEST_BODY)
	}
}

func BenchmarkCreateHoles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// prepare
		b.StopTimer()
		route := "/api/divisions/" + strconv.Itoa(rand.Intn(DIVISION_MAX)+1) + "/holes/"
		data := Map{
			"content": fmt.Sprintf("%v", rand.Uint64()),
			"tag": []Map{
				{"name": "123"},
				{"name": "456"},
			},
		}
		b.StartTimer()

		benchmarkCommon(b, "post", route, REQUEST_BODY, data)
	}
}

func BenchmarkGetHole(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// prepare
		b.StopTimer()
		holeID := rand.Intn(HOLE_MAX) + 1
		url := "/api/holes/" + strconv.Itoa(holeID) + "/"
		b.StartTimer()

		benchmarkCommon(b, "get", url, REQUEST_BODY)
	}
}
