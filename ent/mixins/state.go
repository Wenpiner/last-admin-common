package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type StateMixin struct {
	ent.Schema
}

func (StateMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("state").
			Default(true).
			Optional().
			Comment("状态 / State"),
	}
}
