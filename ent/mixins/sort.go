package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type SortMixin struct {
	mixin.Schema
}

func (SortMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("sort").Default(0).Comment("排序 / Sort"),
	}
}
