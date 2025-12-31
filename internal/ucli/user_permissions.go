package ucli

import (
	"context"
	"fmt"
	"sort"

	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/db"
	"github.com/axelrhd/hagg/internal/user"
	storeUserSqlite "github.com/axelrhd/hagg/internal/user/store_sqlite"
	"github.com/urfave/cli/v3"
)

func userPermissionsCmd() *cli.Command {
	return &cli.Command{
		Name:      "permissions",
		Usage:     "Show effective permissions for users",
		ArgsUsage: "[display-name]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "tree",
				Usage: "Show permission tree (roles -> permissions)",
			},
		},
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

			if c.Bool("tree") {
				return runPermissionsTree(ctx, cfg, userStore, displayName)
			}

			return runPermissionsFlat(ctx, cfg, userStore, displayName)
		},
	}
}

func runPermissionsFlat(
	ctx context.Context,
	cfg *config.Config,
	userStore user.Store,
	displayName string,
) error {

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

		raw, err := enf.GetImplicitPermissionsForUser(sub)
		if err != nil {
			return err
		}

		perms := normalizePerms(raw)
		printYAMLList("permissions", perms, "  ")
		fmt.Println()
	}

	return nil
}

func runPermissionsTree(
	ctx context.Context,
	cfg *config.Config,
	userStore user.Store,
	displayName string,
) error {

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

		fmt.Println("  roles:")

		for _, role := range roles {
			rawPerms, err := enf.GetPermissionsForUser(role)
			if err != nil {
				return err
			}

			perms := normalizePerms(rawPerms)

			// ⬇️ Rolle als Key, Permissions direkt als Liste
			fmt.Printf("    %s:\n", role)
			for _, p := range perms {
				fmt.Printf("      - %s\n", p)
			}
		}

		rawEff, err := enf.GetImplicitPermissionsForUser(sub)
		if err != nil {
			return err
		}

		effective := normalizePerms(rawEff)
		printYAMLList("effective_permissions", effective, "  ")
		fmt.Println()
	}

	return nil
}

func normalizePerms(raw [][]string) []string {
	set := make(map[string]struct{})

	for _, p := range raw {
		// expected: p = [sub, act]
		if len(p) >= 2 {
			act := p[1]
			if act == "*" {
				act = "any"
			}
			set[act] = struct{}{}
		}
	}

	list := make([]string, 0, len(set))
	for k := range set {
		list = append(list, k)
	}
	sort.Strings(list)

	return list
}
