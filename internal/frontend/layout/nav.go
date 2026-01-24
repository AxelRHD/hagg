package layout

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg-lib/view"
	"github.com/axelrhd/hagg/internal/app"
	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func Navbar(ctx *handler.Context, deps app.Deps, oob bool) g.Node {
	isAuthenticated := deps.Auth.IsAuthenticated(ctx.Req)

	return Nav(
		ID("navbar"),
		Class("navbar navbar-expand-lg bg-body-tertiary"),
		g.If(oob,
			hx.SwapOOB("navbar"),
		),

		Div(
			Class("container"),

			// Brand
			A(
				Class("navbar-brand"),
				Href(view.URLString(ctx.Req, "/")),
				Strong(
					Style("font-family: 'Courier New', monospace; letter-spacing: 0.02em;"),
					g.Text("HAGG"),
				),
			),

			// Mobile toggle button
			Button(
				Class("navbar-toggler"),
				Type("button"),
				g.Attr("data-bs-toggle", "collapse"),
				g.Attr("data-bs-target", "#navbarNav"),
				g.Attr("aria-controls", "navbarNav"),
				g.Attr("aria-expanded", "false"),
				g.Attr("aria-label", "Toggle navigation"),
				Span(Class("navbar-toggler-icon")),
			),

			// Collapsible content
			Div(
				Class("collapse navbar-collapse"),
				ID("navbarNav"),

				// Left nav items
				Ul(
					Class("navbar-nav me-auto"),
					Li(
						Class("nav-item"),
						A(
							Class("nav-link"),
							Href(view.URLString(ctx.Req, "/")),
							g.Text("Home"),
						),
					),

					// Show Dashboard link if authenticated
					g.If(isAuthenticated,
						Li(
							Class("nav-item"),
							A(
								Class("nav-link"),
								Href(view.URLString(ctx.Req, "/dashboard")),
								g.Text("Dashboard"),
							),
						),
					),
				),

				// Right nav items
				Ul(
					Class("navbar-nav align-items-center"),

					// Show Login link if not authenticated
					g.If(!isAuthenticated,
						Li(
							Class("nav-item"),
							A(
								Class("nav-link"),
								Href(view.URLString(ctx.Req, "/login")),
								g.Text("Login"),
							),
						),
					),

					// Show Logout button if authenticated
					g.If(isAuthenticated,
						Li(
							Class("nav-item"),
							Button(
								Class("btn btn-outline-secondary btn-sm"),
								hx.Post(view.URLString(ctx.Req, "/htmx/logout")),
								g.Text("Logout"),
							),
						),
					),

					// Theme toggle (always visible)
					Li(
						Class("nav-item ms-2"),
						Button(
							Class("theme-toggle btn btn-link"),
							x.On("click", "theme = theme === 'light' ? 'dark' : 'light'"),

							// Moon icon (shows when light mode)
							g.Raw(`<i class="bi bi-moon-fill" x-show="theme !== 'dark'"></i>`),
							// Sun icon (shows when dark mode)
							g.Raw(`<i class="bi bi-sun-fill" x-show="theme === 'dark'"></i>`),
						),
					),
				),
			),
		),
	)
}
