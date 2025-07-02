package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/servling/servling/pkg/util"
)

// Domain holds the schema definition for the Domain entity.
type Domain struct {
	ent.Schema
}

// Fields of the Domain.
func (Domain) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(util.NewNanoID).
			Unique().
			Immutable(),
		field.String("name").Unique(),
		field.String("certificate").Optional().Nillable(),
		field.String("key").Optional().Nillable(),
		field.String("cloudflare_email").Optional().Nillable(),
		field.String("cloudflare_api_key").Optional().Nillable(),
	}
}

// Edges of the Domain.
func (Domain) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ingresses", Ingress.Type),
	}
}
