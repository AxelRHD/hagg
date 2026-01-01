package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/middleware"
	"github.com/axelrhd/hagg/internal/session"
)

// NewChiServer creates and configures a new Chi router with middleware stack.
//
// Middleware order is important:
//  1. Built-in Chi middleware (RealIP, Compress)
//  2. SCS Session middleware (MUST come first before any middleware that uses sessions)
//  3. Custom middleware (Recovery, Logger, CORS, RateLimit, Secure)
//
// Routes are added via AddRoutes.
//
// Example:
//
//	wrapper := handler.NewWrapper(slog.Default())
//	router := server.NewChiServer(wrapper, deps)
//	http.ListenAndServe(":8080", router)
func NewChiServer(wrapper *handler.Wrapper, deps app.Deps) http.Handler {
	r := chi.NewRouter()

	// Built-in Chi middleware
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Compress(5))

	// SCS Session middleware - MUST come before any middleware that uses sessions!
	r.Use(session.Manager.LoadAndSave)

	// Custom middleware
	r.Use(middleware.Recovery(wrapper))
	r.Use(middleware.Logger(wrapper))
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit)
	r.Use(middleware.Secure)

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Add application routes
	AddRoutes(r, wrapper, deps)

	return r
}
