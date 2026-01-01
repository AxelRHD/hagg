package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/axelrhd/hagg-lib/handler"
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
// Routes should be added to the returned router in routing.go.
//
// Example:
//
//	wrapper := handler.NewWrapper(slog.Default())
//	router := server.NewChiServer(wrapper)
//	http.ListenAndServe(":8080", router)
func NewChiServer(wrapper *handler.Wrapper) http.Handler {
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

	// Routes will be added here (in routing.go migration)
	// For now, add a simple test route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Chi server is running!"))
	})

	return r
}
