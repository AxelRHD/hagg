package hagg

import (
	"github.com/go-chi/chi/v5"

	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/pages/dashboard"
	"github.com/axelrhd/hagg/internal/frontend/pages/home"
	"github.com/axelrhd/hagg/internal/frontend/pages/login"
	"github.com/axelrhd/hagg/internal/middleware"
)

// AddRoutes configures all HTTP routes for the application.
// It registers:
//   - Page routes (full HTML pages): /, /login, /dashboard
//   - HTMX routes (partial HTML): /htmx/login, /htmx/logout
//
// Routes are protected by authentication middleware where appropriate.
func AddRoutes(r chi.Router, wrapper *handler.Wrapper, deps app.Deps) {
	// Public routes
	// Homepage
	r.Get("/", wrapper.Wrap(home.Page(deps)))

	// Login page (GET and POST both render the page)
	// POST is needed for HTMX auto-refresh on auth-changed event
	r.Get("/login", wrapper.Wrap(login.Page(deps)))
	r.Post("/login", wrapper.Wrap(login.Page(deps)))

	// HTMX authentication endpoints
	r.Post("/htmx/login", wrapper.Wrap(login.HxLogin(deps)))
	r.Post("/htmx/logout", wrapper.Wrap(login.HxLogout(deps)))

	// Protected routes (require authentication only)
	// Use RequireAuth for routes that just need a logged-in user
	// Example: r.Use(middleware.RequireAuth(wrapper))

	// Protected routes (require authentication + permission)
	// The dashboard demonstrates Casbin-based permission checks.
	// Users need the "dashboard:view" action assigned to their role.
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission(deps.Auth, deps.Users, deps.Perms, "dashboard:view"))

		r.Get("/dashboard", wrapper.Wrap(dashboard.Page(deps)))
	})
}
