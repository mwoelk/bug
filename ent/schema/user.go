package schema

import (
	"context"
	"entgo.io/bug/ent/hook"
	"entgo.io/bug/ent/privacy"
	"entgo.io/bug/rule"
	"log"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	gen "entgo.io/bug/ent"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age").
			Positive(),
		field.String("name").
			Default("unknown"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.UserFunc(func(ctx context.Context, m *gen.UserMutation) (ent.Value, error) {
					ids, err := m.IDs(ctx)
					if err != nil {
						return nil, err
					}

					log.Printf("--> Hook called with ids %v", ids)

					if len(ids) > 0 {
						log.Fatalf("Hook should not receive any ids due to privacy filters")
					}

					return next.Mutate(ctx, m)
				})
			},
			ent.OpDelete|ent.OpDeleteOne|ent.OpUpdate|ent.OpUpdateOne,
		),
	}
}

// Policy defines the privacy policy of the User.
func (User) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			rule.UserMutationRule(),
		},
		Query: privacy.QueryPolicy{
			// Allow any viewer to read anything.
			privacy.AlwaysAllowRule(),
		},
	}
}
