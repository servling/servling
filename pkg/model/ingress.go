package model

import "github.com/servling/servling/ent"

type CreateIngressInput struct {
	Name       string `json:"name"`
	ServiceId  string `json:"service_id"`
	TargetPort int    `json:"target_port"`
}

type Ingress struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	TargetPort uint16 `json:"target_port"`

	// Relationships from edges
	Domain  *Domain  `json:"domain,omitempty"`
	Service *Service `json:"service,omitempty"`
}

func ingressFromParents(i *ent.Ingress, parentDomain *Domain, parentService *Service) *Ingress {
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
	} else if i.Edges.Domain != nil {
		ingress.Domain = DomainFromEnt(i.Edges.Domain)
	}

	if parentService != nil {
		ingress.Service = parentService
	} else if i.Edges.Service != nil {
		ingress.Service = ServiceFromEnt(i.Edges.Service)
	}

	return ingress
}

func IngressFromEnt(i *ent.Ingress) *Ingress {
	if i == nil {
		return nil
	}
	// Call the helper with nil parents.
	return ingressFromParents(i, nil, nil)
}
