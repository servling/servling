package domain

import (
	"context"

	"dario.lol/gotils/pkg/slice"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/model"
)

//goland:noinspection GoNameStartsWithPackageName
type DomainService struct {
	repository *DomainRepository
}

func NewDomainService(client *ent.Client) *DomainService {
	return &DomainService{repository: NewDomainRepository(client)}
}

func (s *DomainService) Create(ctx context.Context, input model.CreateDomainInput) (*model.Domain, error) {
	dom, err := s.repository.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return model.DomainFromEnt(dom), nil
}

func (s *DomainService) GetByID(ctx context.Context, id string) (*model.Domain, error) {
	dom, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.DomainFromEnt(dom), nil
}

func (s *DomainService) GetByName(ctx context.Context, name string) (*model.Domain, error) {
	dom, err := s.repository.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return model.DomainFromEnt(dom), nil
}

func (s *DomainService) GetAll(ctx context.Context) ([]*model.Domain, error) {
	domains, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return slice.Map(domains, model.DomainFromEnt), err
}

func (s *DomainService) Delete(ctx context.Context, dom *model.Domain) (*model.Domain, error) {
	err := s.repository.Delete(ctx, dom.ID)
	if err != nil {
		return nil, err
	}
	return dom, nil
}
