package domain

import (
	"context"

	"github.com/servling/servling/ent"
	"github.com/servling/servling/ent/domain"
	"github.com/servling/servling/pkg/model"
)

//goland:noinspection GoNameStartsWithPackageName
type DomainRepository struct {
	client *ent.Client
}

func NewDomainRepository(client *ent.Client) *DomainRepository {
	return &DomainRepository{
		client: client,
	}
}

func (r *DomainRepository) GetAll(ctx context.Context) ([]*ent.Domain, error) {
	return r.client.Domain.Query().All(ctx)
}

func (r *DomainRepository) GetByID(ctx context.Context, id string) (*ent.Domain, error) {
	return r.client.Domain.Get(ctx, id)
}

func (r *DomainRepository) GetByName(ctx context.Context, name string) (*ent.Domain, error) {
	return r.client.Domain.Query().Where(domain.Name(name)).Only(ctx)
}

func (r *DomainRepository) GetOrCreateByName(ctx context.Context, name string) (*ent.Domain, error) {
	foundDomain, err := r.GetByName(ctx, name)
	if !ent.IsNotFound(err) {
		foundDomain, err = r.client.Domain.Create().SetName(name).Save(ctx)
		if err != nil {
			return nil, err
		}
		return foundDomain, nil
	}
	if err != nil {
		return nil, err
	}
	return foundDomain, nil
}

func (r *DomainRepository) Create(ctx context.Context, input model.CreateDomainInput) (*ent.Domain, error) {
	return r.client.Domain.Create().
		SetName(input.Name).
		SetNillableCertificate(input.Certificate).
		SetNillableKey(input.Key).
		SetNillableCloudflareEmail(input.CloudflareEmail).
		SetNillableCloudflareAPIKey(input.CloudflareAPIKey).
		Save(ctx)
}

func (r *DomainRepository) Delete(ctx context.Context, id string) error {
	return r.client.Domain.DeleteOneID(id).Exec(ctx)
}
