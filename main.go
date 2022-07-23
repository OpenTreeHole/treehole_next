package main

import (
	"treehole_next/bootstrap"
	"treehole_next/utils"
)

// @title Open Tree Hole
// @version 2.0.0
// @description An Anonymous BBS \n Note: PUT methods are used to PARTLY update, and we don't use PATCH method.

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
	app := bootstrap.Init()
	defer utils.Logger.Sync()
	err := app.Listen("0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
}
