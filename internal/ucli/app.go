package ucli

import (
	"context"

	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "kl-toolbox",
		Usage: "KL Toolbox server and admin CLI",

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
