package model

import "github.com/servling/servling/ent"

type Domain struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Certificate      *string `json:"certificate,omitempty"`
	Key              *string `json:"key,omitempty"`
	CloudflareEmail  *string `json:"cloudflare_email,omitempty"`
	CloudflareAPIKey *string `json:"cloudflare_api_key,omitempty"`

	Ingresses []*Ingress `json:"ingresses,omitempty"`
}

type CreateDomainInput struct {
	Name             string  `json:"name"`
	Certificate      *string `json:"certificate,omitempty"`
	Key              *string `json:"key,omitempty"`
	CloudflareEmail  *string `json:"cloudflare_email,omitempty"`
	CloudflareAPIKey *string `json:"cloudflare_api_key,omitempty"`
}

func DomainFromEnt(d *ent.Domain) *Domain {
	if d == nil {
		return nil
	}
	domain := &Domain{
		ID:               d.ID,
		Name:             d.Name,
		Certificate:      d.Certificate,
		Key:              d.Key,
		CloudflareEmail:  d.CloudflareEmail,
		CloudflareAPIKey: d.CloudflareAPIKey,
	}

	if d.Edges.Ingresses != nil {
		domain.Ingresses = make([]*Ingress, len(d.Edges.Ingresses))
		for i, ingressEnt := range d.Edges.Ingresses {
			domain.Ingresses[i] = ingressFromParents(ingressEnt, domain, nil)
		}
	}

	return domain
}
