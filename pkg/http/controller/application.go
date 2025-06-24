package controller

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/servling/servling/pkg/constants"
	"github.com/servling/servling/pkg/domain/application"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/http/custom_option"
	"github.com/servling/servling/pkg/http/handler"
	"github.com/servling/servling/pkg/types"
)

type ApplicationController struct {
	authService        *auth.AuthService
	applicationService *application.ApplicationService
}

func NewApplicationController(applicationService *application.ApplicationService, authService *auth.AuthService) *ApplicationController {
	return &ApplicationController{
		applicationService: applicationService,
		authService:        authService,
	}
}

func (ac *ApplicationController) Routes(server *fuego.Server) {
	applicationRoutes := fuego.Group(server, "/applications", custom_option.RequirePasetoAuth(ac.authService))

	fuego.Get(applicationRoutes, "/", ac.GetAll, option.OperationID("get-applications"))
	fuego.Post(applicationRoutes, "/", ac.Create, option.OperationID("create-application"))
	fuego.Get(applicationRoutes, "/{id}", ac.Get, option.OperationID("get-application"))
	fuego.Delete(applicationRoutes, "/{id}", ac.Delete, option.OperationID("delete-application"))
	fuego.Post(applicationRoutes, "/{id}/start", ac.Start, option.OperationID("start-application"))
	fuego.Post(applicationRoutes, "/{id}/stop", ac.Stop, option.OperationID("stop-application"))
	fuego.Get(applicationRoutes, "/events", ac.Events, option.OperationID("get-application-events"))
}

func (ac *ApplicationController) Get(c fuego.Context[any, any]) (*types.Application, error) {
	return ac.applicationService.GetByID(c, c.PathParam("id"))
}

func (ac *ApplicationController) GetAll(c fuego.Context[any, any]) ([]types.Application, error) {
	apps, err := ac.applicationService.GetAll(c)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

type CreateApplicationRequest struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Start       bool                       `json:"start"`
	Services    []types.CreateServiceInput `json:"services"`
}

func (ac *ApplicationController) Create(c fuego.Context[CreateApplicationRequest, any]) (*types.Application, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	return ac.applicationService.Create(c, types.CreateApplicationInput{
		Name:        body.Name,
		Description: body.Description,
		Start:       body.Start,
		Services:    body.Services,
	})
}

func (ac *ApplicationController) Delete(c fuego.Context[any, any]) (*types.Application, error) {
	app, err := ac.applicationService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	return ac.applicationService.Delete(c, app)
}

func (ac *ApplicationController) Start(c fuego.Context[any, any]) (*types.Application, error) {
	app, err := ac.applicationService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	if app.Status == types.ServiceStatusStarting {
		return app, nil
	}
	go ac.applicationService.Start(context.Background(), app)
	return app, nil
}

func (ac *ApplicationController) Stop(c fuego.Context[any, any]) (*types.Application, error) {
	app, err := ac.applicationService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	if app.Status == types.ServiceStatusStopping {
		return app, nil
	}
	go ac.applicationService.Stop(context.Background(), app)
	return app, nil
}

func (ac *ApplicationController) Events(c fuego.Context[any, any]) (*types.ApplicationStatusChangedMessage, error) {
	return handler.SSEEventsController[types.ApplicationStatusChangedMessage](c, ac.applicationService.GetPubSub(), constants.TopicApplicationStatusChanged)
}
