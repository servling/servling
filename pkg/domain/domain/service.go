package domain

import "github.com/servling/servling/ent"

//goland:noinspection GoNameStartsWithPackageName
type DomainService struct {
	client *ent.Client
}

func NewDomainService(client *ent.Client) DomainService {
	return DomainService{client: client}
}
