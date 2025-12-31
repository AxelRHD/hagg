package user

import "context"

type Store interface {
	FindByUID(ctx context.Context, uid string) (*User, error)
	FindByDisplayName(ctx context.Context, displayName string) (*User, error)
	CreateUser(ctx context.Context, uid, displayName string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
}
