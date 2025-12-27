package storesqlite

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/axelrhd/hagg/internal/user"
)

func mapSQLError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return user.ErrNotFound
	}

	// optional: UNIQUE constraint
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return user.ErrAlreadyExists
	}

	return err
}
