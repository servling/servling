package server

import (
	"context"
	"os"
	"time"

	"entgo.io/ent/dialect"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/alexdrl/zerowater"
	"github.com/docker/docker/client"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/config"
	"github.com/servling/servling/pkg/deploy"
	"github.com/servling/servling/pkg/deploy/runtime"
	"github.com/servling/servling/pkg/http"
	"github.com/servling/servling/pkg/util"
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
	entClient, err := ent.Open(dialect.Postgres, servlingConfig.Database.ToPostgresDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed opening connection to postgres")
		return
	}
	defer util.CloserOrLog(entClient, "failed closing connection to sqlite")
	if err := entClient.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed creating schema resources")
		return
	}

	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		zerowater.NewZerologLoggerAdapter(log.Logger),
	)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating docker entClient")
		return
	}

	deployManager := deploy.NewDeployManager(runtime.NewDockerRuntime(dockerClient, pubSub), pubSub)

	deployManager.WatchForServiceStatusInfoUpdates(context.Background())

	httpServer := http.NewHttpServer(servlingConfig, entClient, pubSub, deployManager)
	err = httpServer.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed starting http server")
	}
}
