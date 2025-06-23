package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/servling/servling/pkg/util"
)

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

// Fields of the Template.
func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(util.NewNanoID).
			Unique().
			Immutable(),
		field.String("name"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("applications", Application.Type),
	}
}
