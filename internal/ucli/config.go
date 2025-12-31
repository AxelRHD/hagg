package ucli

import (
	"context"

	"github.com/axelrhd/hagg/internal/config"
	"github.com/urfave/cli/v3"
)

func configCmd() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Configuration utilities",
		Action: func(ctx context.Context, _ *cli.Command) error {
			cfg := config.MustLoad()
			cfg.Print()
			return nil
		},
	}
}
