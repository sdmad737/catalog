package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/sdmad737/catalog/backend/internal/data/ent/schema/mixins"
)

// ItemEvent stores a human-readable operational audit trail for an item.
type ItemEvent struct{ ent.Schema }

func (ItemEvent) Mixin() []ent.Mixin { return []ent.Mixin{mixins.BaseMixin{}} }

func (ItemEvent) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").NotEmpty().MaxLen(64),
		field.String("summary").NotEmpty().MaxLen(500),
		field.String("note").Optional().MaxLen(1500),
		field.String("actor_name").Optional().MaxLen(255),
	}
}

func (ItemEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("item", Item.Type).Ref("events").Unique().Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}
