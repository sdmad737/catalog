package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/sdmad737/catalog/backend/internal/data/ent/schema/mixins"
)

// ItemTag is a revocable public identity shared by an NFC tag and its QR fallback.
type ItemTag struct{ ent.Schema }

func (ItemTag) Mixin() []ent.Mixin { return []ent.Mixin{mixins.BaseMixin{}} }

func (ItemTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("public_token").Immutable().NotEmpty().MaxLen(64).Unique(),
		field.Time("revoked_at").Optional().Nillable(),
		field.Time("verified_at").Optional().Nillable(),
		field.Time("last_scanned_at").Optional().Nillable(),
	}
}

func (ItemTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("item", Item.Type).Ref("tags").Unique().Required().
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}

func (ItemTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("revoked_at"),
		index.Edges("item").Unique().Annotations(entsql.IndexWhere("revoked_at IS NULL")),
	}
}
