package model

import (
	"time"

	"github.com/servling/servling/ent"
)

// CreateApplicationInput defines the structure for creating a new application.
type CreateApplicationInput struct {
	Name        string               `json:"name" validate:"required"`
	Description string               `json:"description" validate:"required"`
	Start       bool                 `json:"start" validate:"required"`
	Services    []CreateServiceInput `json:"services" validate:"required"`
}

// CreateServiceInput defines the structure for a service within a new application.
type CreateServiceInput struct {
	Name        string            `json:"name" validate:"required"`
	Image       string            `json:"image" validate:"required"`
	Entrypoint  string            `json:"entrypoint" validate:"required"`
	Environment map[string]string `json:"environment" validate:"required"`
	Ports       map[string]string `json:"ports" validate:"required"`
	Labels      map[string]string `json:"labels" validate:"required"`
}

// The Application represents the structure of an application that is returned from the API.
//
//goland:noinspection GoNameStartsWithPackageName
type Application struct {
	ID          string        `json:"id" validate:"required"`
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Services    []*Service    `json:"services" validate:"required"`
	Status      ServiceStatus `json:"status" validate:"required"  enum:"running,stopped,starting,stopping,error"`
	Error       *string       `json:"error"`
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

// ServiceStatusInfo holds the status information of a container.
type ServiceStatusInfo struct {
	Status ServiceStatus `json:"status" validate:"required" enum:"running,stopped,starting,stopping,error"`
	Error  *string       `json:"error,omitempty"`
}

// ServiceStatusInfoUpdate is used for broadcasting container status updates.
type ServiceStatusInfoUpdate struct {
	ServiceStatusInfo
	ID string `json:"id" validate:"required"`
}

// ApplicationStatusInfoUpdate is used for broadcasting container status updates.
type ApplicationStatusInfoUpdate struct {
	ServiceStatusInfo
	ID string `json:"id" validate:"required"`
}

// Service represents the structure of a service that is returned from the API.
type Service struct {
	ID          string            `json:"id" validate:"required"`
	Name        string            `json:"name" validate:"required"`
	ServiceName string            `json:"serviceName" validate:"required"`
	Image       string            `json:"image" validate:"required"`
	Environment map[string]string `json:"environment" validate:"required"`
	Ports       map[string]string `json:"ports" validate:"required"`
	Labels      map[string]string `json:"labels" validate:"required"`
	Status      ServiceStatus     `json:"status" validate:"required" enum:"running,stopped,starting,stopping,error"`
	Error       *string           `json:"error"`
	Ingresses   []*Ingress        `json:"ingresses" validate:"required"`
	Application *Application      `json:"-" validate:"required"`
	CreatedAt   time.Time         `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time         `json:"updatedAt" validate:"required"`
}

func ApplicationFromEnt(app *ent.Application) *Application {
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
	}

	if app.Edges.Services != nil {
		application.Services = make([]*Service, len(app.Edges.Services))
		for i, serviceEnt := range app.Edges.Services {
			application.Services[i] = serviceFromParent(serviceEnt, application)
		}
	}

	return application
}

func ServiceFromEnt(s *ent.Service) *Service {
	if s == nil {
		return nil
	}
	return serviceFromParent(s, nil)
}

func serviceFromParent(s *ent.Service, parentApp *Application) *Service {
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

	if parentApp != nil {
		service.Application = parentApp
	} else if s.Edges.Application != nil {
		service.Application = ApplicationFromEnt(s.Edges.Application)
	}

	if s.Edges.Ingresses != nil {
		service.Ingresses = make([]*Ingress, len(s.Edges.Ingresses))
		for i, ingressEnt := range s.Edges.Ingresses {
			service.Ingresses[i] = ingressFromParents(ingressEnt, nil, service)
		}
	}

	return service
}
