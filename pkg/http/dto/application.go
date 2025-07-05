package dto

import (
	"time"

	"dario.lol/gotils/pkg/slice"
	"github.com/servling/servling/pkg/model"
)

//goland:noinspection GoSnakeCaseUsage
type CreateApplicationRequest struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Start       bool                       `json:"start"`
	Services    []model.CreateServiceInput `json:"services"`
}

//goland:noinspection GoSnakeCaseUsage
type Application struct {
	ID          string        `json:"id" validate:"required"`
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Status      ServiceStatus `json:"status" validate:"required"  enum:"running,stopped,starting,stopping,error"`
	Error       *string       `json:"error"`
	Services    []*Service    `json:"services"`
	CreatedAt   time.Time     `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time     `json:"updatedAt" validate:"required"`
}

type ServiceStatus string

const (
	ServiceStatusRunning  ServiceStatus = "running"
	ServiceStatusStopped  ServiceStatus = "stopped"
	ServiceStatusStarting ServiceStatus = "starting"
	ServiceStatusStopping ServiceStatus = "stopping"
	ServiceStatusError    ServiceStatus = "error"
)

type Service struct {
	ID            string            `json:"id" validate:"required"`
	Name          string            `json:"name" validate:"required"`
	ServiceName   string            `json:"serviceName" validate:"required"`
	Image         string            `json:"image" validate:"required"`
	Environment   map[string]string `json:"environment" validate:"required"`
	Ports         map[string]string `json:"ports" validate:"required"`
	Labels        map[string]string `json:"labels" validate:"required"`
	Status        ServiceStatus     `json:"status" validate:"required" enum:"running,stopped,starting,stopping,error"`
	Error         *string           `json:"error"`
	Ingresses     []*Ingress        `json:"ingresses" validate:"required"`
	ApplicationID string            `json:"applicationId" validate:"required"`
	CreatedAt     time.Time         `json:"createdAt" validate:"required"`
	UpdatedAt     time.Time         `json:"updatedAt" validate:"required"`
}

func ApplicationFromModel(app *model.Application) *Application {
	if app == nil {
		return nil
	}
	application := &Application{
		ID:        app.ID,
		Name:      app.Name,
		Status:    ServiceStatus(app.Status),
		Error:     app.Error,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
		Services:  slice.Map(app.Services, ServiceFromModel),
	}

	return application
}

func ServiceFromModel(s *model.Service) *Service {
	if s == nil {
		return nil
	}
	return serviceFromParent(s, "")
}

func serviceFromParent(s *model.Service, parentAppID string) *Service {
	if s == nil {
		return nil
	}

	service := &Service{
		ID:          s.ID,
		Name:        s.Name,
		ServiceName: s.ServiceName,
		Image:       s.Image,
		Environment: s.Environment,
		Ports:       s.Ports,
		Labels:      s.Labels,
		Status:      ServiceStatus(s.Status),
		Error:       s.Error,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	if parentAppID != "" {
		service.ApplicationID = parentAppID
	} else if s.Application != nil {
		service.ApplicationID = s.Application.ID
	}

	if s.Ingresses != nil {
		service.Ingresses = make([]*Ingress, len(s.Ingresses))
		for i, ingressEnt := range s.Ingresses {
			service.Ingresses[i] = IngressFromModel(ingressEnt)
		}
	}

	return service
}
