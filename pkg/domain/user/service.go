package user

import (
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/config"
)

//goland:noinspection GoNameStartsWithPackageName
type UserService struct {
	config config.Config
	client ent.Client
}

func NewUserService(config config.Config, client ent.Client) *UserService {
	return &UserService{config: config, client: client}
}
