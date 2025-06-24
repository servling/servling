package server

import (
	"context"
	"os"
	"time"

	"entgo.io/ent/dialect"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/config"
	"github.com/servling/servling/pkg/deploy"
	"github.com/servling/servling/pkg/deploy/runtime"
	"github.com/servling/servling/pkg/http"

	_ "github.com/mattn/go-sqlite3"
)

func initLogger() {
	env := os.Getenv("APP_ENV")

	if env == "production" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().With().Str("scope", "generic").Logger()
		log.Info().Msg("Running in Production Mode")
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Str("scope", "generic").Logger()
		log.Info().Msg("Running in Development Mode")
	}
}

func Run() {
	initLogger()
	servlingConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
		return
	}
	entClient, err := ent.Open(dialect.SQLite, "file:ent.db?_fk=1")
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to sqlite")
		return
	}
	defer entClient.Close()
	if err := entClient.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed creating schema resources")
		return
	}

	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating docker entClient")
		return
	}

	deployManager := deploy.NewDeployManager(runtime.NewDockerRuntime(dockerClient, pubSub))

	httpServer := http.NewHttpServer(servlingConfig, entClient, pubSub, deployManager)
	err = httpServer.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed starting http server")
	}
}
