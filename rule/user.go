package rule

import (
	"context"
	"entgo.io/bug/ent"
	"entgo.io/bug/ent/privacy"
	"entgo.io/bug/ent/user"
	"log"
)

func UserMutationRule() privacy.MutationRule {
	return privacy.UserMutationRuleFunc(func(ctx context.Context, m *ent.UserMutation) error {
		log.Printf("--> Evaluating user policy")

		m.Where(
			user.Age(17),
		)

		return privacy.Skip
	})
}
