package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/sdmad737/catalog/backend/internal/data/ent/schema/mixins"
)

// Person is a borrower or assignee. People are deliberately separate from
// staff accounts: a student can receive an item without being able to sign in.
type Person struct{ ent.Schema }

func (Person) Mixin() []ent.Mixin {
	return []ent.Mixin{mixins.BaseMixin{}, GroupMixin{ref: "people"}}
}

func (Person) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().MaxLen(255),
		field.String("email").Optional().MaxLen(255),
		field.String("identifier").Optional().MaxLen(100),
		field.String("notes").Optional().MaxLen(1000),
		field.Bool("active").Default(true),
	}
}

func (Person) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("checkouts", Checkout.Type).
			Annotations(entsql.Annotation{OnDelete: entsql.Restrict}),
	}
}
