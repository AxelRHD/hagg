package storesqlite

import (
	"context"

	"github.com/axelrhd/hagg/internal/user"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

// Compile-time interface check
var _ user.Store = (*Store)(nil)

func (s *Store) CreateUser(ctx context.Context, uid string) (*user.User, error) {
	q := qCreateUser(uid)

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err // Programmierfehler
	}

	var u user.User
	if err := s.db.GetContext(ctx, &u, sql, args...); err != nil {
		return nil, mapSQLError(err)
	}

	return &u, nil
}

func (s *Store) FindByUID(ctx context.Context, uid string) (*user.User, error) {
	q := qUserByUID(uid)

	sql, args, err := q.ToSql()
	if err != nil {
		// SQL konnte nicht gebaut werden → Programmierfehler
		return nil, err
	}

	var u user.User
	if err := s.db.GetContext(ctx, &u, sql, args...); err != nil {
		return nil, mapSQLError(err)
	}

	return &u, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]*user.User, error) {
	q := qAllUsers()

	sql, args, err := q.ToSql()
	if err != nil {
		// SQL konnte nicht gebaut werden → Programmierfehler
		return nil, err
	}

	var users []*user.User
	if err := s.db.SelectContext(ctx, &users, sql, args...); err != nil {
		return nil, mapSQLError(err)
	}

	return users, nil
}
