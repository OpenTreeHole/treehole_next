package tests

import (
	"github.com/gofiber/fiber/v2"
	"testing"
)

func TestIndex(t *testing.T) {
	testAPI(t, "get", "/", 200, nil, fiber.Map{"message": "hello world"})
	testAPI(t, "get", "/404", 404, nil, fiber.Map{"message": "Cannot GET /404"})
}

func TestDocs(t *testing.T) {
	testCommon(t, "get", "/docs", 302)
	testCommon(t, "get", "/docs/index.html", 200)
}
