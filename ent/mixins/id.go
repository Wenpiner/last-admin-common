package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// ID64Mixin UInt64Mixin implements the uint64 id pattern.
type ID64Mixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (ID64Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Positive().Immutable(),
	}
}

type IDMixin struct {
	mixin.Schema
}

func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Positive().Immutable(),
	}
}

// ID32Mixin UInt32Mixin implements the uint32 id pattern.
type ID32Mixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (ID32Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id").Positive().Immutable(),
	}
}

// UUIDMixin implements the uuid id pattern.
type UUIDMixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (UUIDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(GetUUID).Comment("ID / ID"),
	}
}

// GetUUID returns the uuid of the entity.
func GetUUID() uuid.UUID {
	v7, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return v7
}
