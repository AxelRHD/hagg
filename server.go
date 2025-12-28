package hagg

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/axelrhd/hagg-lib/casbinx"
	libmw "github.com/axelrhd/hagg-lib/middleware"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/auth"
	"github.com/axelrhd/hagg/internal/config"
	"github.com/axelrhd/hagg/internal/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//go:embed static/*
var embeddedFs embed.FS

// StartServer startet den HTTP-Server entweder
// - über TCP (DEV)
// - oder über Unix-Socket (PROD)
func StartServer(cfg *config.Config, usrStore user.Store) {
	router := buildRouter(cfg, usrStore)

	// 🔀 Socket oder TCP?
	if cfg.Server.Socket != "" {
		startUnixSocket(router, cfg.Server.Socket)
		return
	}

	// klassisch (DEV)
	log.Println("Listening on", cfg.Addr())
	if err := router.Run(cfg.Addr()); err != nil {
		log.Fatal(err)
	}
}

// ------------------------------------------------------------
// Router-Build (transport-agnostisch)
// ------------------------------------------------------------

func buildRouter(cfg *config.Config, usrStore user.Store) *gin.Engine {
	r := gin.New()

	// session
	store := cookie.NewStore([]byte(cfg.Session.Secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(cfg.Session.MaxAge.Seconds()),
		HttpOnly: true,
		Secure:   false, // true bei HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// casbin plugin
	enforcer, err := casbinx.NewFileEnforcer(
		cfg.Casbin.ModelPath,
		cfg.Casbin.PolicyPath,
	)
	if err != nil {
		log.Fatal(err)
	}

	// deps
	deps := app.Deps{
		Users:    usrStore,
		Auth:     auth.New(usrStore),
		Enforcer: enforcer,
	}

	// middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions(cfg.Session.CookieName, store))
	r.Use(libmw.BasePath(cfg.Server.BasePath))
	r.Use(libmw.HXTriggers())

	// static files
	staticFs, err := fs.Sub(embeddedFs, "static")
	if err != nil {
		log.Fatal(err)
	}

	rg := r.Group(cfg.Server.BasePath)
	rg.StaticFS("/static", http.FS(staticFs))

	Routing(rg, deps)

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
