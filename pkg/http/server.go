package http

import (
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/cors"
	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/config"
	"github.com/servling/servling/pkg/deploy"
	"github.com/servling/servling/pkg/domain/application"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/http/controller"
)

//goland:noinspection ALL
type HttpServer struct {
	config        *config.Config
	client        *ent.Client
	pubSub        *gochannel.GoChannel
	deployManager *deploy.DeployManager
}

func convertLogLevel(level zerolog.Level) slog.Level {
	var slogLevel slog.Level
	switch level {
	case zerolog.TraceLevel:
		slogLevel = slog.LevelDebug
	case zerolog.DebugLevel:
		slogLevel = slog.LevelDebug
	case zerolog.InfoLevel:
		slogLevel = slog.LevelInfo
	case zerolog.WarnLevel:
		slogLevel = slog.LevelWarn
	case zerolog.ErrorLevel:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo // Default to Info
	}
	return slogLevel
}

func NewHttpServer(config *config.Config, client *ent.Client, pubSub *gochannel.GoChannel, deployManager *deploy.DeployManager) *HttpServer {
	return &HttpServer{
		config:        config,
		client:        client,
		pubSub:        pubSub,
		deployManager: deployManager,
	}
}

func parseEnumTag(tag reflect.StructTag, schema *openapi3.Schema) {
	enumTag, ok := tag.Lookup("enum")
	if !ok {
		return
	}

	enumValues := strings.Split(enumTag, ",")
	if len(enumValues) == 1 && enumValues[0] == "" {
		return
	}

	schema.Enum = make([]interface{}, len(enumValues))
	for i, v := range enumValues {
		schema.Enum[i] = v
	}
}

func (s *HttpServer) Run() error {
	server := fuego.NewServer(
		fuego.WithSecurity(openapi3.SecuritySchemes{
			"PasetoAuth": {
				Value: &openapi3.SecurityScheme{
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "PASETO",
					Description:  "Authenticate requests using a PASETO token in the `Authorization` header.",
				},
			},
		}),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: func(specURL string) http.Handler {
					return fuego.DefaultOpenAPIHandler(specURL)
				},
				JSONFilePath: "./schema/openapi.json",
			}),
			fuego.WithOpenAPIGeneratorSchemaCustomizer(func(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
				parseEnumTag(tag, schema)
				return nil
			}),
		),
		fuego.WithGlobalMiddlewares(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		})),
		fuego.WithLogHandler(slogzerolog.Option{
			Level:  convertLogLevel(log.Logger.GetLevel()),
			Logger: &log.Logger,
		}.NewZerologHandler()),
	)

	authService := auth.NewAuthService(s.config, s.client)
	authController := controller.NewAuthController(authService)
	authController.Routes(server)

	applicationService := application.NewApplicationService(s.config, s.client, s.pubSub, s.deployManager)
	go func() {
		err := applicationService.SubscribeToServiceEvents()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to subscribe to service events")
		}
	}()
	applicationController := controller.NewApplicationController(applicationService, authService)
	applicationController.Routes(server)

	return server.Run()
}
