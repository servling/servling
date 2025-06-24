package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/types"
)

func ProvideAuthContext(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			authHeader := r.Header.Get("Authorization")

			if authHeader != "" {
				log.Debug().Str("Authorization", authHeader).Msg("ProvideAuthContext")
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenString = parts[1]
				} else {
					log.Debug().Str("method", "ProvideAuthContext").Msg("Invalid Authorization header")
					next.ServeHTTP(w, r)
					return
				}
			} else {
				tokenFromQuery := r.URL.Query().Get("token")
				if tokenFromQuery != "" {
					log.Debug().Str("token", tokenFromQuery).Msg("Token found in query parameter for SSE request")
					tokenString = tokenFromQuery
				}
			}

			if tokenString == "" {
				next.ServeHTTP(w, r)
				log.Debug().Str("method", "ProvideAuthContext").Msg("No token found in header or query")
				return
			}

			log.Debug().Str("tokenString", tokenString).Msg("ProvideAuthContext")

			claims, err := authService.VerifyAccessToken(tokenString)
			if err != nil {
				log.Debug().Str("method", "ProvideAuthContext").Err(err).Msg("Invalid token")
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), auth.CtxUserKey{}, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := r.Context().Value(auth.CtxUserKey{}).(*types.AccessTokenPayload)

			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Authentication is required"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
