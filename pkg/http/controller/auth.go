package controller

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/types"
)

type AuthController struct {
	authService *auth.AuthService
}

func NewAuthController(authService *auth.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Routes(server *fuego.Server) {
	authRoutes := fuego.Group(server, "/auth")

	fuego.Post(authRoutes, "/register", ac.Register, option.OperationID("register"))
	fuego.Post(authRoutes, "/login", ac.Login, option.OperationID("login"))
	fuego.Post(authRoutes, "/refresh", ac.Refresh, option.OperationID("refresh"))
	fuego.Post(authRoutes, "/invalidate", ac.Invalidate, option.OperationID("invalidate"))
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ac *AuthController) Register(c fuego.Context[RegisterRequest, any]) (*types.RegisterResult, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	return ac.authService.Register(c, body.Username, body.Password)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ac *AuthController) Login(c fuego.Context[LoginRequest, any]) (*types.LoginResult, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	return ac.authService.Login(c, body.Username, body.Password)
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (ac *AuthController) Refresh(c fuego.Context[RefreshRequest, any]) (*types.RefreshResult, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	return ac.authService.Refresh(c, body.RefreshToken)
}

func (ac *AuthController) Invalidate(c fuego.ContextNoBody) (*types.InvalidateResult, error) {
	err := ac.authService.Invalidate(c)
	if err != nil {
		return nil, err
	}
	return &types.InvalidateResult{
		Ok: true,
	}, nil
}
