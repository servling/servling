package dto

import "github.com/servling/servling/pkg/model"

type ApplicationStatusChangedMessage struct {
	ID     string        `json:"id" validate:"required"`
	Status ServiceStatus `json:"status" validate:"required" enum:"running,stopped,starting,stopping,error"`
	Error  *string       `json:"error,omitempty"`
}

func ApplicationStatusChangedMessageFromModel(m *model.ApplicationStatusChangedMessage) *ApplicationStatusChangedMessage {
	return &ApplicationStatusChangedMessage{
		ID:     m.ID,
		Status: ServiceStatus(m.Status),
		Error:  m.Error,
	}
}

type ServiceStatusChangedMessage struct {
	ID     string        `json:"id" validate:"required"`
	Status ServiceStatus `json:"status" validate:"required" enum:"running,stopped,starting,stopping,error"`
	Error  *string       `json:"error,omitempty"`
}

func ServiceStatusChangedMessageFromModel(m *model.ServiceStatusChangedMessage) *ServiceStatusChangedMessage {
	return &ServiceStatusChangedMessage{
		ID:     m.ID,
		Status: ServiceStatus(m.Status),
		Error:  m.Error,
	}
}
