package login

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg-lib/view"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/layout"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

// Page is the login page handler.
// It renders the login form or logout button depending on authentication status.
func Page(deps app.Deps) handler.HandlerFunc {
	return func(ctx *handler.Context) error {
		user, _ := deps.Auth.CurrentUser(ctx.Req)

		loginURL := view.URLString(ctx.Req, "/htmx/login")
		logoutURL := view.URLString(ctx.Req, "/htmx/logout")

		// Extract username safely
		username := ""
		if user != nil {
			username = user.FullName()
		}

		content := Div(
			// HTMX auto-refresh on auth-changed event
			hx.Post(view.URLString(ctx.Req, "/login")),
			hx.Trigger("auth-changed from:body"),
			hx.Target("#page"),
			hx.Select("#page"),
			hx.Swap("outerHTML"),

			// Centering with Bootstrap utilities
			Class("d-flex align-items-center justify-content-center p-3"),
			Style("min-height: 80vh"),

			g.If(!deps.Auth.IsAuthenticated(ctx.Req),
				LoginForm(loginURL),
			),

			g.If(deps.Auth.IsAuthenticated(ctx.Req) && username != "",
				LogoutButton(logoutURL, username),
			),
		)

		return ctx.Render(layout.Page(ctx, deps, content))
	}
}
