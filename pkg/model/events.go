package model

type ApplicationStatusChangedMessage struct {
	ID     string        `json:"id" validate:"required"`
	Status ServiceStatus `json:"status" validate:"required" enum:"running,stopped,starting,stopping,error"`
	Error  *string       `json:"error,omitempty"`
}

type ServiceStatusChangedMessage struct {
	ID     string        `json:"id" validate:"required"`
	Status ServiceStatus `json:"status" validate:"required" enum:"running,stopped,starting,stopping,error"`
	Error  *string       `json:"error,omitempty"`
}
