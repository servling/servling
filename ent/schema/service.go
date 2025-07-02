package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/servling/servling/pkg/util"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(util.NewNanoID).
			Unique().
			Immutable(),
		field.String("name").Unique(),
		field.String("service_name").Unique(),
		field.String("image"),
		field.JSON("ports", map[string]string{}).
			Optional(),
		field.JSON("environment", map[string]string{}).
			Optional(),
		field.String("entrypoint").
			Optional(),
		field.JSON("labels", map[string]string{}).
			Optional(),
		field.String("status").Default("stopped"),
		field.String("error").Optional().Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("application", Application.Type).
			Ref("services").
			Unique(),
		edge.To("ingresses", Ingress.Type),
	}
}
