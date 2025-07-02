package auth

import (
	"context"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
	"dario.lol/gotils/pkg/encoding"
	"dario.lol/gotils/pkg/hash"
	"github.com/go-fuego/fuego"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/config"
	"github.com/servling/servling/pkg/domain/user"
	"github.com/servling/servling/pkg/model"
)

type CtxUserKey struct{}

//goland:noinspection GoNameStartsWithPackageName
type AuthService struct {
	config         *config.Config
	userRepository *user.UserRepository
}

func NewAuthService(config *config.Config, client *ent.Client) *AuthService {
	return &AuthService{config: config, userRepository: user.NewUserRepository(client)}
}

func (s *AuthService) Register(ctx context.Context, username, password string) (*model.RegisterResult, error) {
	hashedPassword, err := hash.Argon2idStringToString(password)

	if err != nil {
		return nil, err
	}

	databaseUser, err := s.userRepository.Create(ctx, model.CreateUserInput{
		Username:       username,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	resultUser := model.UserFromEnt(databaseUser)

	accessToken, err := s.GenerateAccessToken(*resultUser)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	refreshToken, err := s.GenerateRefreshToken(*resultUser)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	return &model.RegisterResult{
		User:                  *resultUser,
		AccessToken:           accessToken.Token,
		AccessTokenExpiresAt:  accessToken.ExpiresAt,
		RefreshToken:          refreshToken.Token,
		RefreshTokenExpiresAt: refreshToken.ExpiresAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*model.LoginResult, error) {
	databaseUser, err := s.userRepository.GetByName(ctx, username)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	verified, err := hash.VerifyArgon2idString(databaseUser.Password, password)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	if !verified {
		return nil, errors.New("invalid password")
	}
	resultUser := model.UserFromEnt(databaseUser)

	accessToken, err := s.GenerateAccessToken(*resultUser)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	refreshToken, err := s.GenerateRefreshToken(*resultUser)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	return &model.LoginResult{
		User:                  *resultUser,
		AccessToken:           accessToken.Token,
		AccessTokenExpiresAt:  accessToken.ExpiresAt,
		RefreshToken:          refreshToken.Token,
		RefreshTokenExpiresAt: refreshToken.ExpiresAt,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.RefreshResult, error) {
	tokenPayload, err := s.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	databaseUser, err := s.userRepository.GetByID(ctx, tokenPayload.ID)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	if tokenPayload.TokenVersion != databaseUser.TokenVersion {
		return nil, errors.New("invalid refresh token")
	}
	resultUser := model.UserFromEnt(databaseUser)
	accessToken, err := s.GenerateAccessToken(*resultUser)
	if err != nil {
		return nil, fuego.ForbiddenError{Detail: err.Error()}
	}
	return &model.RefreshResult{
		AccessToken:          accessToken.Token,
		AccessTokenExpiresAt: accessToken.ExpiresAt,
	}, nil
}

func (s *AuthService) Invalidate(ctx context.Context) error {
	currentUser, err := s.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	return s.userRepository.IncrementTokenVersion(ctx, currentUser.ID)
}

func (s *AuthService) GetUserFromContext(ctx context.Context) (*model.User, error) {
	tokenPayload, ok := ctx.Value(CtxUserKey{}).(*model.AccessTokenPayload)
	if !ok {
		return nil, errors.New("invalid token")
	}
	databaseUser, err := s.userRepository.GetByID(ctx, tokenPayload.ID)
	if err != nil {
		return nil, err
	}
	return model.UserFromEnt(databaseUser), nil
}

func (s *AuthService) GenerateAccessToken(user model.User) (*model.TokenResult, error) {
	token := paseto.NewToken()
	expiresAt := time.Now().Add(s.config.Security.Token.AccessTokenDuration)

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(expiresAt)

	token.SetSubject(user.ID)
	token.SetString("id", user.ID)
	token.SetString("username", user.Name)
	token.SetTime("createdAt", user.CreatedAt)
	token.SetTime("updatedAt", user.UpdatedAt)

	secretKey, err := paseto.V4SymmetricKeyFromBytes(s.config.Security.Token.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	return &model.TokenResult{
		Token:     token.V4Encrypt(secretKey, nil),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) GenerateRefreshToken(user model.User) (*model.TokenResult, error) {
	token := paseto.NewToken()
	expiresAt := time.Now().Add(s.config.Security.Token.RefreshTokenDuration)

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(expiresAt)

	token.SetSubject(user.ID)
	token.SetString("id", user.ID)
	token.SetString("username", user.Name)
	err := token.Set("tokenVersion", user.TokenVersion)
	if err != nil {
		return nil, err
	}

	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(s.config.Security.Token.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	return &model.TokenResult{
		Token:     token.V4Sign(privateKey, nil),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) VerifyRefreshToken(tokenString string) (*model.RefreshTokenPayload, error) {
	key, err := paseto.NewV4AsymmetricPublicKeyFromBytes(s.config.Security.Token.RefreshTokenPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := paseto.NewParser().ParseV4Public(key, tokenString, nil)
	if err != nil {
		return nil, err
	}
	payload, err := encoding.UnmarshalJSON[*model.RefreshTokenPayload](token.ClaimsJSON())
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (s *AuthService) VerifyAccessToken(tokenString string) (*model.AccessTokenPayload, error) {
	key, err := paseto.V4SymmetricKeyFromBytes(s.config.Security.Token.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}
	token, err := paseto.NewParser().ParseV4Local(key, tokenString, nil)
	if err != nil {
		return nil, err
	}
	payload, err := encoding.UnmarshalJSON[*model.AccessTokenPayload](token.ClaimsJSON())
	if err != nil {
		return nil, err
	}

	return payload, nil
}
