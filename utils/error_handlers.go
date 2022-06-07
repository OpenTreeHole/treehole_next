package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MyErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	code := 500
	message := err.Error()

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		code = 404
	}

	return ctx.Status(code).JSON(fiber.Map{"message": message})
}
