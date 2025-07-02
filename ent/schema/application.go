package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/servling/servling/pkg/util"
)

// Application holds the schema definition for the Application entity.
type Application struct {
	ent.Schema
}

// Fields of the Application.
func (Application) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(util.NewNanoID).
			Unique().
			Immutable(),
		field.String("name").Unique(),
		field.String("description"),
		field.String("image_url").
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

// Edges of the Application.
func (Application) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("services", Service.Type),

		edge.From("template", Template.Type).
			Ref("applications").
			Unique(),
	}
}
