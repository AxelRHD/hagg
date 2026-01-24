package ucli

import (
	"context"

	"github.com/axelrhd/hagg/internal/version"
	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:    "hagg",
		Usage:   "HAGG Stack server and admin CLI",
		Version: version.Version,

		EnableShellCompletion: true,

		// DEFAULT
		Action: func(_ context.Context, _ *cli.Command) error {
			return serve()
		},

		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start the HTTP server",
				Action: func(_ context.Context, _ *cli.Command) error {
					return serve()
				},
			},
			configCmd(),
			userCmd(),
		},
	}
}
