package controller

import (
	"dario.lol/gotils/pkg/slice"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/domain/domain"
	"github.com/servling/servling/pkg/http/custom_option"
	"github.com/servling/servling/pkg/http/dto"
)

type DomainController struct {
	authService   *auth.AuthService
	domainService *domain.DomainService
}

func NewDomainController(domainService *domain.DomainService, authService *auth.AuthService) *DomainController {
	return &DomainController{
		domainService: domainService,
		authService:   authService,
	}
}

func (ac *DomainController) Routes(server *fuego.Server) {
	applicationRoutes := fuego.Group(server, "/domains", custom_option.RequirePasetoAuth(ac.authService))

	fuego.Get(applicationRoutes, "/", ac.GetAll, option.OperationID("get-applications"))
	fuego.Post(applicationRoutes, "/", ac.Create, option.OperationID("create-application"))
	fuego.Get(applicationRoutes, "/{id}", ac.Get, option.OperationID("get-application"))
	fuego.Delete(applicationRoutes, "/{id}", ac.Delete, option.OperationID("delete-application"))
}

func (ac *DomainController) Get(c fuego.Context[any, any]) (*dto.Domain, error) {
	app, err := ac.domainService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	return dto.DomainFromModel(app), nil
}

func (ac *DomainController) GetAll(c fuego.Context[any, any]) ([]*dto.Domain, error) {
	apps, err := ac.domainService.GetAll(c)
	if err != nil {
		return nil, err
	}
	return slice.Map(apps, dto.DomainFromModel), nil
}

func (ac *DomainController) Create(c fuego.Context[dto.CreateDomainRequest, any]) (*dto.Domain, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	app, err := ac.domainService.Create(c, body.ToInput())
	if err != nil {
		return nil, err
	}
	return dto.DomainFromModel(app), nil
}

func (ac *DomainController) Delete(c fuego.Context[any, any]) (*dto.Domain, error) {
	app, err := ac.domainService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	deletedApp, err := ac.domainService.Delete(c, app)
	if err != nil {
		return nil, err
	}
	return dto.DomainFromModel(deletedApp), nil
}
