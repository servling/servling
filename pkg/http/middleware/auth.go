package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/model"
	"github.com/servling/servling/pkg/util"
)

func ProvideAuthContext(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string
			authHeader := r.Header.Get("Authorization")

			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenString = parts[1]
				} else {
					next.ServeHTTP(w, r)
					return
				}
			} else {
				tokenFromQuery := r.URL.Query().Get("token")
				if tokenFromQuery != "" {
					tokenString = tokenFromQuery
				}
			}

			if tokenString == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := authService.VerifyAccessToken(tokenString)
			if err != nil {
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
			_, ok := r.Context().Value(auth.CtxUserKey{}).(*model.AccessTokenPayload)

			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				util.MustOrLog(w.Write([]byte(`{"error": "Authentication is required"}`)))("Error writing response")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
