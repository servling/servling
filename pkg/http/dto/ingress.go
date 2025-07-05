package dto

import (
	"github.com/servling/servling/pkg/model"
)

type Ingress struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	TargetPort uint16 `json:"target_port"`

	Domain    *Domain `json:"domain,omitempty"`
	ServiceID *string `json:"service_id,omitempty"`
}

func ingressFromParents(i *model.Ingress, parentDomain *Domain, parentServiceId *string) *Ingress {
	if i == nil {
		return nil
	}

	ingress := &Ingress{
		ID:         i.ID,
		Name:       i.Name,
		TargetPort: i.TargetPort,
	}

	if parentDomain != nil {
		ingress.Domain = parentDomain
	} else if i.Domain != nil {
		ingress.Domain = DomainFromModel(i.Domain)
	}

	if parentServiceId != nil {
		ingress.ServiceID = parentServiceId
	} else if i.Service != nil {
		serviceID := i.Service.ID
		ingress.ServiceID = &serviceID
	}

	return ingress
}

func IngressFromModel(i *model.Ingress) *Ingress {
	if i == nil {
		return nil
	}
	return ingressFromParents(i, nil, nil)
}
