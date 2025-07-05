package dto

import "github.com/servling/servling/pkg/model"

type Domain struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Certificate      *string  `json:"certificate,omitempty"`
	Key              *string  `json:"key,omitempty"`
	CloudflareEmail  *string  `json:"cloudflare_email,omitempty"`
	CloudflareAPIKey *string  `json:"cloudflare_api_key,omitempty"`
	IngressIDs       []string `json:"ingress_ids"`
}

func DomainFromModel(d *model.Domain) *Domain {
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

	return domain
}

type CreateDomainRequest struct {
	Name             string  `json:"name" validate:"required"`
	Certificate      *string `json:"certificate,omitempty"`
	Key              *string `json:"key,omitempty"`
	CloudflareEmail  *string `json:"cloudflare_email,omitempty"`
	CloudflareAPIKey *string `json:"cloudflare_api_key,omitempty"`
}

func (req CreateDomainRequest) ToInput() model.CreateDomainInput {
	return model.CreateDomainInput{
		Name:             req.Name,
		Certificate:      req.Certificate,
		Key:              req.Key,
		CloudflareEmail:  req.CloudflareEmail,
		CloudflareAPIKey: req.CloudflareAPIKey,
	}
}
