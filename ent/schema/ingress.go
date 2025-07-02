package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/servling/servling/pkg/util"
)

// Ingress holds the schema definition for the Ingress entity.
type Ingress struct {
	ent.Schema
}

// Fields of the Ingress.
func (Ingress) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(util.NewNanoID).
			Unique().
			Immutable(),
		field.String("name").Unique(),
		field.Uint16("target_port"),
	}
}

// Edges of the Ingress.
func (Ingress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("domain", Domain.Type).
			Ref("ingresses").
			Unique(),
		edge.From("service", Service.Type).
			Ref("ingresses").
			Unique(),
	}
}
