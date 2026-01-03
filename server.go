package hagg

import (
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/axelrhd/hagg-lib/casbinx"
	"github.com/axelrhd/hagg-lib/handler"
	libmw "github.com/axelrhd/hagg-lib/middleware"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/auth"
	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/middleware"
	"github.com/axelrhd/hagg/internal/session"
	"github.com/axelrhd/hagg/internal/user"
)

// StartServer initializes and starts the HTTP server.
// It supports two modes:
//   - TCP mode (development): Uses host:port from config
//   - Unix socket mode (production): Uses socket path from config
//
// The server will block until an error occurs or the process is terminated.
func StartServer(cfg *config.Config, usrStore user.Store) {
	// Initialize SCS session manager
	if err := session.Init(cfg.Session.DBPath); err != nil {
		log.Fatal("failed to init sessions:", err)
	}

	router := buildRouter(cfg, usrStore)

	// Socket or TCP?
	if cfg.Server.Socket != "" {
		startUnixSocket(router, cfg.Server.Socket)
		return
	}

	// TCP mode (dev)
	log.Println("Listening on", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), router); err != nil {
		log.Fatal("Server error:", err)
	}
}

// buildRouter constructs the Chi router with all middleware, dependencies, and routes.
func buildRouter(cfg *config.Config, usrStore user.Store) http.Handler {
	// Create logger
	logger := slog.Default()

	// Create handler wrapper
	wrapper := handler.NewWrapper(logger)

	// Casbin enforcer
	enforcer, err := casbinx.NewFileEnforcer(
		cfg.Casbin.ModelPath,
		cfg.Casbin.PolicyPath,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Dependencies
	deps := app.Deps{
		Users:    usrStore,
		Auth:     auth.New(usrStore),
		Enforcer: enforcer,
		Perms:    casbinx.NewPerm(enforcer),
	}

	// Create Chi router
	r := chi.NewRouter()

	// Built-in Chi middleware
	r.Use(chimw.RealIP)
	r.Use(chimw.Compress(5))

	// SCS Session middleware - MUST come before any middleware that uses sessions!
	r.Use(session.Manager.LoadAndSave)

	// Custom middleware
	r.Use(middleware.Recovery(wrapper))
	r.Use(middleware.Logger(wrapper))
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit)
	r.Use(libmw.Secure)

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Add application routes
	AddRoutes(r, wrapper, deps)

	return r
}

// ------------------------------------------------------------
// Unix-Socket Start
// ------------------------------------------------------------

func startUnixSocket(handler http.Handler, socketName string) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		log.Fatal("XDG_RUNTIME_DIR not set (required for unix socket)")
	}

	socketPath := filepath.Join(runtimeDir, socketName)

	// alten Socket entfernen
	_ = os.Remove(socketPath)

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on unix socket:", socketPath)

	if err := http.Serve(l, handler); err != nil {
		log.Fatal(err)
	}
}
