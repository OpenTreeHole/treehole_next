package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	testAPI(t, "get", "/", 302, nil)
	testAPI(t, "get", "/api", 200, nil)
	web404 := testAPI(t, "get", "/404", 404, nil)
	assert.EqualValues(t, "Cannot GET /404", web404["message"])
}

func TestDocs(t *testing.T) {
	testCommon(t, "get", "/docs", 302)
	testCommon(t, "get", "/docs/index.html", 200)
}
