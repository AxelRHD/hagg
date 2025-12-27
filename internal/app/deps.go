package app

import (
	"github.com/axelrhd/hagg/internal/auth"
	"github.com/axelrhd/hagg/internal/user"

	"github.com/casbin/casbin/v2"
)

type Deps struct {
	Users user.Store
	Auth  *auth.Auth

	// Authorization (RBAC / ABAC)
	Enforcer *casbin.Enforcer
}
