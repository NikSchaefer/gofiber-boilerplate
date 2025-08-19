package schema

import (
	"fmt"
	"math/rand"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Session struct {
	ent.Schema
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.Time("expires").
			Immutable().
			Default(GetTokenExpiration),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Unique().
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}

type OTP struct {
	ent.Schema
}

func (OTP) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (OTP) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty().
			Immutable().
			MaxLen(255).
			DefaultFunc(func() string {
				// Generate a 6-digit OTP code
				return fmt.Sprintf("%06d", rand.Intn(1000000))
			}),
		field.Enum("type").
			Values("login", "password_reset").
			Default("login"),
		field.Bool("used").
			Default(false),
		field.Time("expires_at").
			Immutable().
			Default(GetTokenExpiration),
	}
}

func (OTP) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("otps").
			Unique().
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
