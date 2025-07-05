package dto

import (
	"time"

	"github.com/servling/servling/pkg/model"
)

type LoginResult struct {
	User                  User      `json:"user" validate:"required"`
	AccessToken           string    `json:"accessToken" validate:"required"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt" validate:"required"`
	RefreshToken          string    `json:"refreshToken" validate:"required"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt" validate:"required"`
}

func LoginResultFromModel(result *model.LoginResult) *LoginResult {
	return &LoginResult{
		User:                  *UserFromModel(&result.User),
		AccessToken:           result.AccessToken,
		AccessTokenExpiresAt:  result.AccessTokenExpiresAt,
		RefreshToken:          result.RefreshToken,
		RefreshTokenExpiresAt: result.RefreshTokenExpiresAt,
	}
}

type RegisterResult struct {
	User                  User      `json:"user" validate:"required"`
	AccessToken           string    `json:"accessToken" validate:"required"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt" validate:"required"`
	RefreshToken          string    `json:"refreshToken" validate:"required"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt" validate:"required"`
}

func RegisterResultFromModel(result *model.RegisterResult) *RegisterResult {
	return &RegisterResult{
		User:                  *UserFromModel(&result.User),
		AccessToken:           result.AccessToken,
		AccessTokenExpiresAt:  result.AccessTokenExpiresAt,
		RefreshToken:          result.RefreshToken,
		RefreshTokenExpiresAt: result.RefreshTokenExpiresAt,
	}
}

type RefreshResult struct {
	AccessToken          string    `json:"accessToken" validate:"required"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt" validate:"required"`
}

func RefreshResultFromModel(result *model.RefreshResult) *RefreshResult {
	return &RefreshResult{
		AccessToken:          result.AccessToken,
		AccessTokenExpiresAt: result.AccessTokenExpiresAt,
	}
}

type InvalidateResult struct {
	Ok bool `json:"ok" validate:"required"`
}

func InvalidateResultFromModel(result *model.InvalidateResult) *InvalidateResult {
	return &InvalidateResult{
		Ok: result.Ok,
	}
}
