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
func Page(ctx *handler.Context, deps app.Deps) error {
	user, _ := deps.Auth.CurrentUser(ctx.Req)

	loginURL := view.URLStringChi(ctx.Req, "/htmx/login")
	logoutURL := view.URLStringChi(ctx.Req, "/htmx/logout")

	// Extract username safely
	username := ""
	if user != nil {
		username = user.FullName()
	}

	content := Div(
		// HTMX auto-refresh on auth-changed event
		hx.Post(view.URLStringChi(ctx.Req, "/")),
		hx.Trigger("auth-changed from:body"),
		hx.Target("#page"),
		hx.Select("#page"),
		hx.Swap("outerHTML"),

		// Centering
		Class("flex items-center justify-center pa3"),
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
