package user

import (
	"fmt"

	"github.com/axelrhd/litetime"
)

type User struct {
	ID          int64         `db:"id"`
	UID         string        `db:"uid"`          // login secret (never shown, never used for authz)
	DisplayName string        `db:"display_name"` // stable, human-readable identifier (Casbin subject)
	FirstName   string        `db:"first_name"`
	LastName    string        `db:"last_name"`
	CreatedAt   litetime.Time `db:"created_at"`
	UpdatedAt   litetime.Time `db:"updated_at"`
}

// -----------------------------------------------------------------------------
// Presentation helpers
// -----------------------------------------------------------------------------

// FullName returns a human-friendly name for UI usage.
// Falls back to DisplayName if first/last name are empty.
func (u User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.DisplayName
	}

	if u.FirstName == "" {
		return u.LastName
	}
	if u.LastName == "" {
		return u.FirstName
	}

	return u.LastName + ", " + u.FirstName
}

// -----------------------------------------------------------------------------
// Authorization helpers
// -----------------------------------------------------------------------------

// Subject returns the Casbin subject identifier for this user.
func (u User) Subject() string {
	return u.DisplayName
}

// -----------------------------------------------------------------------------
// String representations (debug / logging)
// -----------------------------------------------------------------------------

func (u User) String() string {
	return fmt.Sprintf(
		"User(id=%d uid='%s' display_name='%s' first_name='%s' last_name='%s' created_at='%s' updated_at='%s')",
		u.ID,
		u.UID,
		u.DisplayName,
		u.FirstName,
		u.LastName,
		u.CreatedAt,
		u.UpdatedAt,
	)
}

func (u User) GoString() string {
	return fmt.Sprintf(
		"User{id=%d uid=%q display_name=%q first_name=%q last_name=%q created_at=%q updated_at=%q}",
		u.ID,
		u.UID,
		u.DisplayName,
		u.FirstName,
		u.LastName,
		u.CreatedAt,
		u.UpdatedAt,
	)
}
