package server

import (
	"context"
	"log"

	"entgo.io/ent/dialect"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/docker/docker/client"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/config"
	"github.com/servling/servling/pkg/deploy"
	"github.com/servling/servling/pkg/deploy/runtime"
	"github.com/servling/servling/pkg/http"

	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	servlingConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	entClient, err := ent.Open(dialect.SQLite, "file:ent.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
		return
	}
	defer entClient.Close()
	if err := entClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
		return
	}

	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("failed creating docker entClient: %v", err)
		return
	}

	deployManager := deploy.NewDeployManager(runtime.NewDockerRuntime(dockerClient, pubSub))

	httpServer := http.NewHttpServer(servlingConfig, entClient, pubSub, deployManager)
	err = httpServer.Run()
	if err != nil {
		log.Fatal(err)
	}
}
