package model

type ApplicationStatusChangedMessage struct {
	ID     string        `json:"id"`
	Status ServiceStatus `json:"status"`
	Error  *string       `json:"error,omitempty"`
}

type ServiceStatusChangedMessage struct {
	ID     string        `json:"id"`
	Status ServiceStatus `json:"status"`
	Error  *string       `json:"error,omitempty"`
}
