package ucli

import (
	"context"
	"fmt"

	"github.com/axelrhd/kl-toolbox/internal/user"
)

func resolveSubjects(
	ctx context.Context,
	userStore user.Store,
	displayName string,
) ([]string, error) {

	// expliziter Name angegeben
	if displayName != "" {
		u, err := userStore.FindByDisplayName(ctx, displayName)
		if err != nil {
			if err == user.ErrNotFound {
				fmt.Printf("user %q not found\n", displayName)
				return nil, nil // ⬅️ KEIN Fehler, KEINE Subjects
			}
			return nil, err
		}

		return []string{u.DisplayName}, nil
	}

	// kein Name → alle User
	users, err := userStore.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	subs := make([]string, 0, len(users))
	for _, u := range users {
		subs = append(subs, u.DisplayName)
	}

	return subs, nil
}
