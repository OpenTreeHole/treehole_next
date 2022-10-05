package benchmarks

import (
	"math/rand"
	"strconv"
	"testing"
	. "treehole_next/models"
)

func BenchmarkListFloorsInAHole(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		route := "/api/holes/" + strconv.Itoa(rand.Intn(HOLE_MAX)+1) + "/floors/"
		b.StartTimer()

		benchmarkCommon(b, "get", route, REQUEST_BODY)
	}
}

func BenchmarkGetFloor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		route := "/api/floors/" + strconv.Itoa(rand.Intn(FLOOR_MAX)+1) + "/"
		b.StartTimer()

		benchmarkCommon(b, "get", route, REQUEST_BODY)
	}
}

func BenchmarkCreateFloor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		route := "/api/holes/" + strconv.Itoa(rand.Intn(HOLE_MAX)+1) + "/floors/"
		data := Map{
			"content":  strconv.Itoa(rand.Int()),
			"reply_to": 0,
		}
		b.StartTimer()

		benchmarkCommon(b, "post", route, REQUEST_BODY, data)
	}
}
