package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	app := bootstrap.Init()
	go func() {
		err := app.Listen("0.0.0.0:8000")
		if err != nil {
			log.Fatal(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	err := app.Shutdown()
	if err != nil {
		log.Fatal(err)
	}

	err = utils.Logger.Sync()
	if err != nil {
		log.Fatal(err)
	}

}
