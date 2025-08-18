package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email")
			.Unique()
			.NotEmpty()
			.MaxLen(255)
			.Immutable(),
		field.Bool("email_verified")
			.Default(false),
		field.String("phone_number").
			Optional().
			Unique().
			MaxLen(255),
		field.Bool("phone_number_verified").
			Default(false),
	}
}

func (User) Edges() []ent.Edge {
	return nil
}

type Account struct {
	ent.Schema
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values("password", "google", "apple"),
		field.Bytes("password_hash").
			Optional().
			MaxLen(255).
			Sensitive(),
		field.String("provider_id").
			Optional().
			MaxLen(255).
			Unique().
			Sensitive(),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("accounts").
			Unique(),
	}
}

// Profile stores the user profile information
type Profile struct {
	ent.Schema
}

func (Profile) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(255),
		field.String("avatar_url").
			Optional().
			MaxLen(2048).
			Comment("URL to the avatar image"),
		field.String("avatar_key").
			Optional().
			MaxLen(255).
			Comment("S3 object key if avatar is uploaded, empty if external URL"),
		field.Time("birthday").
			Optional(),
	}
}

func (Profile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("profile").
			Unique().
			Required(),
	}
}
