package user

import "context"

type Store interface {
	FindByUID(ctx context.Context, uid string) (*User, error)
	CreateUser(ctx context.Context, uid string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
}
