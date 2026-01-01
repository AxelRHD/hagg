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
		g.If(oob,
			hx.SwapOOB("navbar"),
		),

		Ul(
			Li(
				A(
					Href(view.URLStringChi(ctx.Req, "/")),
					Strong(
						Style("font-family: 'Courier Next', Courier, monospace; font-size: 1.5rem; text-transform: uppercase; letter-spacing: 0em"),
						g.Text("HAGG"),
					),
				),
			),
		),

		Ul(
			Li(
				A(
					Href(view.URLStringChi(ctx.Req, "/")),
					g.Text("Home"),
				),
			),

			// Show Login link if not authenticated
			g.If(!isAuthenticated,
				Li(
					A(
						Href(view.URLStringChi(ctx.Req, "/login")),
						g.Text("Login"),
					),
				),
			),

			// Show Dashboard link if authenticated
			g.If(isAuthenticated,
				Li(
					A(
						Href(view.URLStringChi(ctx.Req, "/dashboard")),
						g.Text("Dashboard"),
					),
				),
			),

			// Show Logout button if authenticated
			g.If(isAuthenticated,
				Li(
					Button(
						x.Bind("class", "theme === 'dark' ? 'contrast' : 'secondary'"),
						hx.Post(view.URLStringChi(ctx.Req, "/htmx/logout")),
						g.Text("Logout"),
					),
				),
			),

			// Theme toggle (always visible)
			Li(
				Button(
					Class("theme-toggle"),
					x.On("click", "theme = theme === 'light' ? 'dark' : 'light'"),

					// Moon icon (shows when light mode or not set)
					g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-lucide="icon" x-show="theme !== 'dark'"><path d="M12 3a6 6 0 009 9 9 9 0 11-9-9z"/></svg>`),
					// Sun icon (shows when dark mode)
					g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" data-lucide="icon" x-show="theme === 'dark'"><circle cx="12" cy="12" r="4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="m4.93 4.93 1.41 1.41"/><path d="m17.66 17.66 1.41 1.41"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="m6.34 17.66-1.41 1.41"/><path d="m19.07 4.93-1.41 1.41"/></svg>`),
				),
			),
		),
	)
}
