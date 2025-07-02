package application

import (
	"context"

	"github.com/pkg/errors"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/ent/application"
	"github.com/servling/servling/ent/service"
	"github.com/servling/servling/pkg/model"
	"github.com/servling/servling/pkg/util"
)

//goland:noinspection GoNameStartsWithPackageName
type ApplicationRepository struct {
	client *ent.Client
}

func NewApplicationRepository(client *ent.Client) *ApplicationRepository {
	return &ApplicationRepository{client: client}
}

func (r *ApplicationRepository) GetAll(ctx context.Context) ([]*ent.Application, error) {
	return r.client.Application.Query().WithServices().All(ctx)
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id string) (*ent.Application, error) {
	return r.client.Application.Query().Where(application.ID(id)).WithServices().Only(ctx)
}

func (r *ApplicationRepository) GetByIDWithIngresses(ctx context.Context, id string) (*ent.Application, error) {
	return r.client.Application.Query().Where(application.ID(id)).WithServices(func(query *ent.ServiceQuery) {
		query.WithIngresses()
	}).Only(ctx)
}

func (r *ApplicationRepository) Delete(ctx context.Context, id string) error {
	return r.client.Application.DeleteOneID(id).Exec(ctx)
}

func (r *ApplicationRepository) GetServiceWithApplicationServices(ctx context.Context, id string) (*ent.Service, error) {
	return r.client.Service.Query().Where(service.IDEQ(id)).WithApplication(func(query *ent.ApplicationQuery) {
		query.WithServices()
	}).Only(ctx)
}

func (r *ApplicationRepository) GetAllServices(ctx context.Context) ([]*ent.Service, error) {
	return r.client.Service.Query().All(ctx)
}

func (r *ApplicationRepository) UpdateServiceStatus(ctx context.Context, id string, info model.ServiceStatusInfo) error {
	return r.client.Service.Update().Where(service.IDEQ(id)).SetStatus(string(info.Status)).SetNillableError(info.Error).Exec(ctx)
}

func (r *ApplicationRepository) UpdateApplicationStatus(ctx context.Context, id string, info model.ServiceStatusInfo) error {
	return r.client.Application.Update().Where(application.IDEQ(id)).SetStatus(string(info.Status)).SetNillableError(info.Error).Exec(ctx)
}

func (r *ApplicationRepository) CreateService(ctx context.Context, applicationName string, start bool, input model.CreateServiceInput) (*ent.Service, error) {
	status := model.ServiceStatusStopped
	if start {
		status = model.ServiceStatusStarting
	}
	return r.client.Service.Create().
		SetName(input.Name).
		SetServiceName(util.NormalizeContainerName(applicationName) + "-" + util.NormalizeContainerName(input.Name)).
		SetImage(input.Image).
		SetEntrypoint(input.Entrypoint).
		SetEnvironment(input.Environment).
		SetPorts(input.Ports).
		SetLabels(input.Labels).
		SetStatus(string(status)).
		Save(ctx)
}

func (r *ApplicationRepository) Create(ctx context.Context, input model.CreateApplicationInput) (*ent.Application, error) {
	if len(input.Services) <= 0 {
		return nil, errors.New("no services to create")
	}
	services := make([]*ent.Service, len(input.Services)-1)

	for _, srv := range input.Services {
		srv, err := r.CreateService(ctx, input.Name, input.Start, srv)
		if err != nil {
			return nil, err
		}
		services = append(services, srv)
	}

	if len(input.Services) != len(services) {
		return nil, errors.New("not all services could be created")
	}

	status := model.ServiceStatusStopped
	if input.Start {
		status = model.ServiceStatusStarting
	}

	app, err := r.client.Application.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetStatus(string(status)).
		AddServices(services...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	app.Edges.Services = services
	return app, nil
}
