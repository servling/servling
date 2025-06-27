package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/servling/servling/cli"
	"github.com/servling/servling/pkg/util"
)

func main() {
	util.StartTimer()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go cli.Execute()
	<-shutdown
	log.Info().Msgf("Finished after %s", util.StopTimer())
}
