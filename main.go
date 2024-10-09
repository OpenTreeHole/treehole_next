package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"treehole_next/bootstrap"
)

//	@title			Open Tree Hole
//	@version		2.1.0
//	@description	An Anonymous BBS \n Note: PUT methods are used to PARTLY update, and we don't use PATCH method.

//	@contact.name	Maintainer Ke Chen
//	@contact.email	dev@danta.tech

//	@license.name	Apache 2.0
//	@license.url	https://www.apache.org/licenses/LICENSE-2.0.html

//	@host
//	@BasePath	/api

func main() {
	app, cancel := bootstrap.Init()
	go func() {
		err := app.Listen("0.0.0.0:8000")
		if err != nil {
			log.Fatal().Err(err).Msg("app listen failed")
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		log.Err(err).Msg("error shutdown app")
	}
	// stop tasks
	cancel()
}
