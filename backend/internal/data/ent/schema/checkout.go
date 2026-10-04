package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/sdmad737/catalog/backend/internal/data/ent/schema/mixins"
)

// Checkout records an item's assignment to a borrower and is retained after
// return so the organization has a reliable operational history.
type Checkout struct{ ent.Schema }

func (Checkout) Mixin() []ent.Mixin {
	return []ent.Mixin{mixins.BaseMixin{}, GroupMixin{ref: "checkouts"}}
}

func (Checkout) Fields() []ent.Field {
	return []ent.Field{
		field.Time("checked_out_at").Default(time.Now),
		field.Time("due_at").Optional().Nillable(),
		field.Time("returned_at").Optional().Nillable(),
		field.Enum("condition_out").Default("unknown").Values("unknown", "excellent", "good", "fair", "poor", "damaged"),
		field.Enum("condition_in").Default("unknown").Values("unknown", "excellent", "good", "fair", "poor", "damaged"),
		field.String("checkout_note").Optional().MaxLen(1000),
		field.String("return_note").Optional().MaxLen(1000),
		field.String("checked_out_by").Optional().MaxLen(255),
		field.String("returned_by").Optional().MaxLen(255),
	}
}

func (Checkout) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("item", Item.Type).Ref("checkouts").Unique().Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.From("person", Person.Type).Ref("checkouts").Unique().Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Restrict}),
	}
}
