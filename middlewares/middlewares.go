package middlewares

import "github.com/gofiber/fiber/v2"

func RegisterMiddlewares(app *fiber.App) {
	app.Use(Logger)
}
