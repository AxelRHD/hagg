package ucli

import (
	"context"
	"fmt"

	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/db"
	storeUserSqlite "github.com/axelrhd/hagg/internal/user/store_sqlite"
	"github.com/urfave/cli/v3"
)

func userRolesCmd() *cli.Command {
	return &cli.Command{
		Name:      "roles",
		Usage:     "Show roles for users",
		ArgsUsage: "[display-name]",
		Action: func(ctx context.Context, c *cli.Command) error {

			var displayName string
			if c.Args().Len() > 0 {
				displayName = c.Args().First()
			}

			cfg := config.MustLoad()

			dbx, err := db.OpenSQLite(cfg.Database.SQLite.Path)
			if err != nil {
				return err
			}
			defer dbx.Close()

			userStore := storeUserSqlite.New(dbx)

			enf, err := loadEnforcer(cfg)
			if err != nil {
				return err
			}

			subjects, err := resolveSubjects(ctx, userStore, displayName)
			if err != nil {
				return err
			}

			for _, sub := range subjects {
				fmt.Printf("%s:\n", sub)

				roles, err := enf.GetRolesForUser(sub)
				if err != nil {
					return err
				}

				for _, role := range roles {
					fmt.Printf("  - %s\n", role)
				}
				fmt.Println()
			}

			return nil
		},
	}
}
