package types

import (
	"time"
)

// LoginResult represents the data returned upon a successful login.
type LoginResult struct {
	User                  User      `json:"user" validate:"required"`
	AccessToken           string    `json:"accessToken" validate:"required"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt" validate:"required"`
	RefreshToken          string    `json:"refreshToken" validate:"required"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt" validate:"required"`
}

// RegisterResult represents the data returned upon a successful registration.
type RegisterResult struct {
	User                  User      `json:"user" validate:"required"`
	AccessToken           string    `json:"accessToken" validate:"required"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt" validate:"required"`
	RefreshToken          string    `json:"refreshToken" validate:"required"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt" validate:"required"`
}

// RefreshResult represents the data returned upon a successful token refresh.
type RefreshResult struct {
	AccessToken          string    `json:"accessToken" validate:"required"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt" validate:"required"`
}

// InvalidateResult represents the response for a successful token invalidation.
type InvalidateResult struct {
	Ok bool `json:"ok" validate:"required"`
}

// AccessTokenPayload defines the PASETO payload for an access token.
type AccessTokenPayload struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RefreshTokenPayload defines the PASETO payload for a refresh token.
type RefreshTokenPayload struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	TokenVersion int    `json:"tokenVersion"`
}

// TokenResult is a generic structure for returning a token and its expiration.
type TokenResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
