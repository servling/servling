package user

import (
	"context"

	"github.com/servling/servling/ent"
	"github.com/servling/servling/ent/user"
	"github.com/servling/servling/pkg/types"
)

//goland:noinspection GoNameStartsWithPackageName
type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, input types.CreateUserInput) (*ent.User, error) {
	return r.client.User.Create().
		SetName(input.Username).
		SetPassword(input.HashedPassword).
		Save(ctx)
}

func (r *UserRepository) GetByName(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().Where(user.NameEQ(username)).First(ctx)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*ent.User, error) {
	return r.client.User.Query().Where(user.IDEQ(id)).First(ctx)
}

func (r *UserRepository) IncrementTokenVersion(ctx context.Context, id string) error {
	return r.client.User.Update().Where(user.IDEQ(id)).AddTokenVersion(1).Exec(ctx)
}
