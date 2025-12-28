package ucli

import (
	kltoolbox "github.com/axelrhd/kl-toolbox"
	"github.com/axelrhd/kl-toolbox/internal/config"
	"github.com/axelrhd/kl-toolbox/internal/db"
	storeHuSqlite "github.com/axelrhd/kl-toolbox/internal/hu/store_sqlite"
	storeUserSqlite "github.com/axelrhd/kl-toolbox/internal/user/store_sqlite"
)

func serve() error {
	cfg := config.MustLoad()

	dbx, err := db.OpenSQLite(cfg.Database.SQLite.Path)
	if err != nil {
		return err
	}
	defer dbx.Close()

	userStore := storeUserSqlite.New(dbx)
	huStore := storeHuSqlite.New(dbx)

	kltoolbox.StartServer(cfg, userStore, huStore)
	return nil
}
