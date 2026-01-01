package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/pages/login"
)

// AddLoginRoutes adds login routes to the Chi router.
// This includes the login page (GET/POST /) and HTMX endpoints for login/logout.
func AddLoginRoutes(r chi.Router, wrapper *handler.Wrapper, deps app.Deps) {
	// Login page (GET and POST both render the page)
	// POST is needed for HTMX auto-refresh on auth-changed event
	r.Get("/", wrapper.Wrap(func(ctx *handler.Context) error {
		return login.Page(ctx, deps)
	}))

	r.Post("/", wrapper.Wrap(func(ctx *handler.Context) error {
		return login.Page(ctx, deps)
	}))

	// HTMX endpoints
	r.Post("/htmx/login", wrapper.Wrap(login.HxLogin(deps)))
	r.Post("/htmx/logout", wrapper.Wrap(login.HxLogout(deps)))
}
