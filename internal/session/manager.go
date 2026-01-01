package session

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	_ "github.com/mattn/go-sqlite3"
)

// Manager is the global session manager instance.
// It's initialized once during application startup via Init().
var Manager *scs.SessionManager

// Init initializes the global session manager with SQLite persistent storage.
// This must be called before starting the server.
//
// The session manager is configured with:
//   - 24 hour session lifetime
//   - HttpOnly cookies (prevents XSS attacks)
//   - SameSite=Lax (CSRF protection)
//   - SQLite backend for persistence across restarts
//   - Shared database with app (default: ./db.sqlite3)
//
// Example:
//
//	if err := session.Init("./db.sqlite3"); err != nil {
//	    log.Fatal("failed to init sessions", "error", err)
//	}
func Init(dbPath string) error {
	// Ensure directory exists (if DB is in subdirectory)
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Open SQLite database for session storage
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Create session manager
	Manager = scs.New()
	Manager.Lifetime = 24 * time.Hour
	Manager.Cookie.Name = "hagg_session"
	Manager.Cookie.HttpOnly = true
	Manager.Cookie.SameSite = http.SameSiteLaxMode

	// Persistent storage - sessions survive server restarts
	store := sqlite3store.New(db)
	Manager.Store = store

	return nil
}
