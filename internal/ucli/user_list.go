package ucli

import (
	"context"
	"os"

	"github.com/axelrhd/kl-toolbox/internal/config"
	"github.com/axelrhd/kl-toolbox/internal/db"
	storeUserSqlite "github.com/axelrhd/kl-toolbox/internal/user/store_sqlite"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"
)

func userListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List users",
		Action: func(ctx context.Context, _ *cli.Command) error {

			cfg := config.MustLoad()

			dbx, err := db.OpenSQLite(cfg.Database.SQLite.Path)
			if err != nil {
				return err
			}
			defer dbx.Close()

			store := storeUserSqlite.New(dbx)

			users, err := store.ListUsers(ctx)
			if err != nil {
				return err
			}

			t := table.New(
				"ID",
				"UID",
				"DISPLAY NAME",
				"FULL NAME",
			)
			t.WithWriter(os.Stdout)

			for _, u := range users {
				t.AddRow(
					u.ID,
					u.UID,
					u.DisplayName,
					u.FullName(),
				)
			}

			t.Print()
			return nil
		},
	}
}
