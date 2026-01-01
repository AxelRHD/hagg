package dashboard

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/layout"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Page renders the protected dashboard page.
// This page is only accessible to authenticated users.
func Page(ctx *handler.Context, deps app.Deps) error {
	user, _ := deps.Auth.CurrentUser(ctx.Req)

	// Get username safely
	username := "User"
	if user != nil {
		username = user.FullName()
	}

	content := Div(
		Class("container"),
		Style("padding: 2rem 0;"),

		// Welcome header
		Header(
			Style("margin-bottom: 2rem;"),
			H1(
				g.Text("Dashboard"),
			),
			P(
				Style("color: var(--muted-color);"),
				g.Text("Welcome back, "),
				Strong(g.Text(username)),
				g.Text("!"),
			),
		),

		// Dashboard content
		Div(
			Class("grid"),

			// User info card
			Article(
				H3(g.Text("👤 User Information")),
				P(
					Strong(g.Text("Name: ")),
					g.Text(username),
				),
				g.If(user != nil,
					P(
						Strong(g.Text("UID: ")),
						Code(g.Text(user.UID)),
					),
				),
				P(
					Strong(g.Text("Status: ")),
					Span(
						Style("color: var(--success-color);"),
						g.Text("✓ Authenticated"),
					),
				),
			),

			// Session info card
			Article(
				H3(g.Text("🔐 Session")),
				P(g.Text("Your session is managed server-side with SCS (Session Cookie Store).")),
				P(
					Strong(g.Text("Storage: ")),
					g.Text("SQLite (persistent across restarts)"),
				),
				P(
					Strong(g.Text("Security: ")),
					g.Text("HTTP-only cookies, secure in production"),
				),
			),
		),

		// Features demonstration
		Article(
			Style("margin-top: 2rem;"),
			H2(g.Text("Protected Features")),
			P(
				g.Text("This page demonstrates a protected route in the HAGG boilerplate. "),
				g.Text("Only authenticated users can access this content."),
			),

			Details(
				Style("margin-top: 1rem;"),
				Summary(g.Text("How does authentication work?")),
				Div(
					Style("padding: 1rem;"),
					P(g.Text("The HAGG boilerplate uses Chi middleware for route protection:")),
					Ul(
						Li(
							Code(g.Text("middleware.RequireAuth()")),
							g.Text(" checks session for user ID"),
						),
						Li(g.Text("Unauthenticated users are redirected to /login")),
						Li(g.Text("Sessions persist across server restarts (SQLite backend)")),
						Li(g.Text("Logout clears the session server-side")),
					),
				),
			),

			Details(
				Style("margin-top: 1rem;"),
				Summary(g.Text("What can you build with HAGG?")),
				Div(
					Style("padding: 1rem;"),
					Ul(
						Li(g.Text("Multi-tenant SaaS applications")),
						Li(g.Text("Internal admin dashboards")),
						Li(g.Text("CRUD applications with real-time updates")),
						Li(g.Text("API-first applications with hypermedia UIs")),
						Li(g.Text("Prototypes and MVPs with minimal dependencies")),
					),
				),
			),
		),

		// Call to action
		Article(
			Style("margin-top: 2rem; text-align: center;"),
			H3(g.Text("Ready to build?")),
			P(
				g.Text("Check out the source code to see how this dashboard is built. "),
				g.Text("All UI is rendered server-side with gomponents."),
			),
			Div(
				Style("margin-top: 1.5rem;"),
				A(
					Href("https://github.com/axelrhd/hagg"),
					Role("button"),
					Class("outline"),
					Target("_blank"),
					g.Text("View Source on GitHub"),
				),
			),
		),
	)

	return ctx.Render(layout.Page(ctx, deps, content))
}
