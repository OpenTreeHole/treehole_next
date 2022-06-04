package main

import (
	"github.com/gofiber/fiber/v2"
	"treehole_next/apis"
	"treehole_next/config"
	"treehole_next/middlewares"
	"treehole_next/utils"
)

// @title Tree Hole
// @version 2.0.0
// @description A anonymous bbs

// @contact.name Maintainer Shi Yue
// @contact.email hasbai@fduhole.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.InitConfig()
	initDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.MyErrorHandler,
	})
	middlewares.RegisterMiddlewares(app)
	apis.RegisterRoutes(app)

	err := app.Listen("0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
}
