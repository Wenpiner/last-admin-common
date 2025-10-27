package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type TimestampMixin struct {
	ent.Schema
}

func (TimestampMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间 / Creation time"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间 / Update time"),
	}
}
