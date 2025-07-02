package ingress

import (
	"context"

	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/domain/domain"
	"github.com/servling/servling/pkg/model"
	"golang.org/x/net/publicsuffix"
)

//goland:noinspection GoNameStartsWithPackageName
type IngressService struct {
	repository       *IngressRepository
	domainRepository *domain.DomainRepository
}

func NewIngressService(client *ent.Client) *IngressService {
	return &IngressService{
		repository:       NewIngressRepository(client),
		domainRepository: domain.NewDomainRepository(client),
	}
}

func (s *IngressService) Create(ctx context.Context, input model.CreateIngressInput) (*model.Ingress, error) {
	eTLD, err := publicsuffix.EffectiveTLDPlusOne(input.Name)
	if err != nil {
		return nil, err
	}
	foundDomain, err := s.domainRepository.GetOrCreateByName(ctx, eTLD)
	if err != nil {
		return nil, err
	}
	ingress, err := s.repository.CreateIngress(ctx, CreateDBIngressInput{
		Name:       input.Name,
		DomainId:   foundDomain.ID,
		ServiceId:  input.ServiceId,
		TargetPort: input.TargetPort,
	})
	if err != nil {
		return nil, err
	}
	return model.IngressFromEnt(ingress), nil
}
