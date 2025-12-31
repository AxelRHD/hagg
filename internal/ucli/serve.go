package ucli

import (
	"github.com/axelrhd/hagg"
	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/db"
	storeUserSqlite "github.com/axelrhd/hagg/internal/user/store_sqlite"
)

func serve() error {
	cfg := config.MustLoad()

	dbx, err := db.OpenSQLite(cfg.Database.SQLite.Path)
	if err != nil {
		return err
	}
	defer dbx.Close()

	userStore := storeUserSqlite.New(dbx)

	hagg.StartServer(cfg, userStore)
	return nil
}
