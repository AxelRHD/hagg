package ucli

import (
	"context"
	"fmt"

	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/db"
	storeUserSqlite "github.com/axelrhd/hagg/internal/user/store_sqlite"
	"github.com/urfave/cli/v3"
)

func userCmd() *cli.Command {
	return &cli.Command{
		Name:  "user",
		Usage: "User management",
		Commands: []*cli.Command{
			userCreateCmd(),
			userListCmd(),
			userRolesCmd(),
			userPermissionsCmd(),
		},
	}
}

func userCreateCmd() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new user",
		Action: func(ctx context.Context, _ *cli.Command) error {

			cfg := config.MustLoad()

			dbx, err := db.OpenSQLite(cfg.Database.SQLite.Path)
			if err != nil {
				return err
			}
			defer dbx.Close()

			input, err := promptCreateUser()
			if err != nil {
				return err
			}

			store := storeUserSqlite.New(dbx)

			u, err := store.CreateUser(ctx, input.UID, input.DisplayName)
			if err != nil {
				return err
			}

			fmt.Printf(
				"✔ user created: id=%d display_name=%s\n",
				u.ID,
				u.DisplayName,
			)

			return nil
		},
	}
}
