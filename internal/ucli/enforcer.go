package ucli

import (
	"github.com/axelrhd/hagg-lib/casbinx"
	"github.com/axelrhd/hagg/internal/config"
	"github.com/casbin/casbin/v2"
)

func loadEnforcer(cfg *config.Config) (*casbin.Enforcer, error) {
	return casbinx.NewFileEnforcer(
		cfg.Casbin.ModelPath,
		cfg.Casbin.PolicyPath,
	)
}
