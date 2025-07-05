package controller

import (
	"context"

	"dario.lol/gotils/pkg/slice"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/servling/servling/pkg/constants"
	"github.com/servling/servling/pkg/domain/application"
	"github.com/servling/servling/pkg/domain/auth"
	"github.com/servling/servling/pkg/http/custom_option"
	"github.com/servling/servling/pkg/http/dto"
	"github.com/servling/servling/pkg/http/handler"
	"github.com/servling/servling/pkg/model"
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

func (ac *ApplicationController) Get(c fuego.Context[any, any]) (*dto.Application, error) {
	app, err := ac.applicationService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	return dto.ApplicationFromModel(app), nil
}

func (ac *ApplicationController) GetAll(c fuego.Context[any, any]) ([]*dto.Application, error) {
	apps, err := ac.applicationService.GetAll(c)
	if err != nil {
		return nil, err
	}
	return slice.Map(apps, dto.ApplicationFromModel), nil
}

func (ac *ApplicationController) Create(c fuego.Context[dto.CreateApplicationRequest, any]) (*dto.Application, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}
	app, err := ac.applicationService.Create(c, model.CreateApplicationInput{
		Name:        body.Name,
		Description: body.Description,
		Start:       body.Start,
		Services:    body.Services,
	})
	if err != nil {
		return nil, err
	}
	return dto.ApplicationFromModel(app), nil
}

func (ac *ApplicationController) Delete(c fuego.Context[any, any]) (*dto.Application, error) {
	app, err := ac.applicationService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	deletedApp, err := ac.applicationService.Delete(c, app)
	if err != nil {
		return nil, err
	}
	return dto.ApplicationFromModel(deletedApp), nil
}

func (ac *ApplicationController) Start(c fuego.Context[any, any]) (*dto.Application, error) {
	app, err := ac.applicationService.GetByIDWithIngresses(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	if app.Status == model.ServiceStatusStarting {
		return dto.ApplicationFromModel(app), nil
	}
	go ac.applicationService.Start(context.Background(), app)
	return dto.ApplicationFromModel(app), nil
}

func (ac *ApplicationController) Stop(c fuego.Context[any, any]) (*dto.Application, error) {
	app, err := ac.applicationService.GetByID(c, c.PathParam("id"))
	if err != nil {
		return nil, err
	}
	if app.Status == model.ServiceStatusStopping {
		return dto.ApplicationFromModel(app), nil
	}
	go ac.applicationService.Stop(context.Background(), app)
	return dto.ApplicationFromModel(app), nil
}

func (ac *ApplicationController) Events(c fuego.Context[any, any]) (*dto.ApplicationStatusChangedMessage, error) {
	return handler.SSEEventsController[dto.ApplicationStatusChangedMessage](c, ac.applicationService.GetPubSub(), constants.TopicApplicationStatusChanged)
}
