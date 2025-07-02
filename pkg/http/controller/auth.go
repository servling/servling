package controller

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/model"
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
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ac *AuthController) Register(c fuego.Context[RegisterRequest, any]) (*model.RegisterResult, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	return ac.authService.Register(c, body.Username, body.Password)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ac *AuthController) Login(c fuego.Context[LoginRequest, any]) (*model.LoginResult, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	return ac.authService.Login(c, body.Username, body.Password)
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func (ac *AuthController) Refresh(c fuego.Context[RefreshRequest, any]) (*model.RefreshResult, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	return ac.authService.Refresh(c, body.RefreshToken)
}

func (ac *AuthController) Invalidate(c fuego.ContextNoBody) (*model.InvalidateResult, error) {
	err := ac.authService.Invalidate(c)
	if err != nil {
		return nil, err
	}
	return &model.InvalidateResult{
		Ok: true,
	}, nil
}
