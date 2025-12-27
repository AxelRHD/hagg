package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/axelrhd/hagg"
	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/db"
	"github.com/axelrhd/hagg/internal/user"
	storeUserSqlite "github.com/axelrhd/hagg/internal/user/store_sqlite"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	printConfig := flag.Bool("config", false, "print configuration")
	newUser := flag.String("new-user", "", "create new user with given uid")
	flag.Parse()

	cfg := config.MustLoad()

	if *printConfig {
		cfg.Print()
		os.Exit(0)
	}

	// --- Composition Root ---
	dbx, err := db.OpenSQLite(cfg.Database.SQLite.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	userStore := storeUserSqlite.New(dbx)

	// --- CLI mode ---
	if *newUser != "" {
		handleCreateUser(userStore, *newUser)
		os.Exit(0)
	}

	// --- Server mode ---
	hacc.StartServer(cfg, userStore)
}

func handleCreateUser(store user.Store, uid string) {
	ctx := context.Background()

	u, err := store.CreateUser(ctx, uid)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"created user: id=%d uid=%s created_at=%s",
		u.ID,
		u.UID,
		u.CreatedAt.Format("2006-01-02 15:04:05"),
	)
}
