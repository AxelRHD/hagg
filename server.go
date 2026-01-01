package hagg

import (
	"embed"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/axelrhd/hagg-lib/casbinx"
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/auth"
	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/server"
	"github.com/axelrhd/hagg/internal/session"
	"github.com/axelrhd/hagg/internal/user"
)

//go:embed static/*
var embeddedFs embed.FS

// StartServer startet den HTTP-Server entweder
// - über TCP (DEV)
// - oder über Unix-Socket (PROD)
func StartServer(cfg *config.Config, usrStore user.Store) {
	// Initialize SCS session manager
	if err := session.Init(cfg.Session.DBPath); err != nil {
		log.Fatal("failed to init sessions:", err)
	}

	// Build deps
	enforcer, err := casbinx.NewFileEnforcer(
		cfg.Casbin.ModelPath,
		cfg.Casbin.PolicyPath,
	)
	if err != nil {
		log.Fatal(err)
	}

	deps := app.Deps{
		Users:    usrStore,
		Auth:     auth.New(usrStore),
		Enforcer: enforcer,
	}

	// Build Chi router
	wrapper := handler.NewWrapper(slog.Default())
	chiRouter := server.NewChiServer(wrapper, deps)

	// 🔀 Socket oder TCP?
	if cfg.Server.Socket != "" {
		startUnixSocket(chiRouter, cfg.Server.Socket)
		return
	}

	// Start Chi server on configured port
	log.Println("Chi server listening on", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), chiRouter); err != nil {
		log.Fatal("Server error:", err)
	}
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
