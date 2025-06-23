package types

import (
	"time"

	"dario.lol/gotils/pkg/slice"
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
	Services    []Service     `json:"services" validate:"required"`
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
	CreatedAt   time.Time         `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time         `json:"updatedAt" validate:"required"`
}

func ServiceFromEnt(service *ent.Service) *Service {
	return &Service{
		ID:          service.ID,
		Name:        service.Name,
		ServiceName: service.ServiceName,
		Image:       service.Image,
		Environment: service.Environment,
		Ports:       service.Ports,
		Labels:      service.Labels,
		Status:      ServiceStatus(service.Status),
		Error:       service.Error,
		CreatedAt:   service.CreatedAt,
		UpdatedAt:   service.UpdatedAt,
	}
}

func ApplicationFromEnt(app *ent.Application) *Application {
	return &Application{
		ID:        app.ID,
		Name:      app.Name,
		Services:  slice.FromPtr(slice.Map(app.Edges.Services, ServiceFromEnt)),
		Status:    ServiceStatus(app.Status),
		Error:     app.Error,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}
}
