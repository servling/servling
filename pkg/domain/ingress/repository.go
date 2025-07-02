package ingress

import (
	"context"

	"github.com/servling/servling/ent"
)

//goland:noinspection GoNameStartsWithPackageName
type IngressRepository struct {
	client *ent.Client
}

func NewIngressRepository(client *ent.Client) *IngressRepository {
	return &IngressRepository{
		client: client,
	}
}

type CreateDBIngressInput struct {
	Name       string `json:"name"`
	DomainId   string `json:"domain_id"`
	ServiceId  string `json:"service_id"`
	TargetPort int    `json:"target_port"`
}

func (r *IngressRepository) CreateIngress(ctx context.Context, input CreateDBIngressInput) (*ent.Ingress, error) {
	return r.client.Ingress.Create().
		SetName(input.Name).
		SetDomainID(input.DomainId).
		SetServiceID(input.ServiceId).
		Save(ctx)
}
