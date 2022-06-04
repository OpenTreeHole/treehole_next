package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MyErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := 500
	message := "Internal Server Error"

	// Retrieve the custom status code if it's a fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = 404
		message = err.Error()
	}

	if err != nil {
		return ctx.Status(code).JSON(fiber.Map{"message": message})
	}

	return nil
}
