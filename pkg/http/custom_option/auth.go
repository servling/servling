package custom_option

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/http/middleware"
)

func WithPasetoAuth(authService *auth.AuthService) func(route *fuego.BaseRoute) {
	return func(route *fuego.BaseRoute) {
		option.Middleware(middleware.ProvideAuthContext(authService))(route)
		option.Security(openapi3.SecurityRequirement{
			"PasetoAuth": []string{},
		})(route)
	}
}

func RequirePasetoAuth(authService *auth.AuthService) func(route *fuego.BaseRoute) {
	return func(route *fuego.BaseRoute) {
		option.Middleware(middleware.ProvideAuthContext(authService), middleware.RequireAuth())(route)
		option.Security(openapi3.SecurityRequirement{
			"PasetoAuth": []string{},
		})(route)
	}
}
