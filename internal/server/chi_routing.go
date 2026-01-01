package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/pages/dashboard"
	"github.com/axelrhd/hagg/internal/frontend/pages/home"
	"github.com/axelrhd/hagg/internal/frontend/pages/login"
	"github.com/axelrhd/hagg/internal/middleware"
)

// AddRoutes adds all application routes to the Chi router.
func AddRoutes(r chi.Router, wrapper *handler.Wrapper, deps app.Deps) {
	// Public routes
	// Homepage
	r.Get("/", wrapper.Wrap(func(ctx *handler.Context) error {
		return home.Page(ctx, deps)
	}))

	// Login page (GET and POST both render the page)
	// POST is needed for HTMX auto-refresh on auth-changed event
	r.Get("/login", wrapper.Wrap(func(ctx *handler.Context) error {
		return login.Page(ctx, deps)
	}))

	r.Post("/login", wrapper.Wrap(func(ctx *handler.Context) error {
		return login.Page(ctx, deps)
	}))

	// HTMX authentication endpoints
	r.Post("/htmx/login", wrapper.Wrap(login.HxLogin(deps)))
	r.Post("/htmx/logout", wrapper.Wrap(login.HxLogout(deps)))

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(wrapper))

		r.Get("/dashboard", wrapper.Wrap(func(ctx *handler.Context) error {
			return dashboard.Page(ctx, deps)
		}))
	})
}
