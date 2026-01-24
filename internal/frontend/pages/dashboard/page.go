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
func Page(deps app.Deps) handler.HandlerFunc {
	return func(ctx *handler.Context) error {
		user, _ := deps.Auth.CurrentUser(ctx.Req)

		// Get username safely
		username := "User"
		if user != nil {
			username = user.FullName()
		}

		content := Div(
			Class("container py-4"),

			// Welcome header
			Header(
				Class("mb-4"),
				H1(
					g.Text("Dashboard"),
				),
				P(
					Class("text-body-secondary"),
					g.Text("Welcome back, "),
					Strong(g.Text(username)),
					g.Text("!"),
				),
			),

			// Dashboard content
			Div(
				Class("row g-4 mb-4"),

				// User info card
				Div(
					Class("col-md-6"),
					Article(
						Class("card h-100 p-4"),
						H3(Class("h5"), g.Text("👤 User Information")),
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
							Class("mb-0"),
							Strong(g.Text("Status: ")),
							Span(
								Class("text-success"),
								g.Text("✓ Authenticated"),
							),
						),
					),
				),

				// Session info card
				Div(
					Class("col-md-6"),
					Article(
						Class("card h-100 p-4"),
						H3(Class("h5"), g.Text("🔐 Session")),
						P(g.Text("Your session is managed server-side with SCS (Session Cookie Store).")),
						P(
							Strong(g.Text("Storage: ")),
							g.Text("SQLite (persistent across restarts)"),
						),
						P(
							Class("mb-0"),
							Strong(g.Text("Security: ")),
							g.Text("HTTP-only cookies, secure in production"),
						),
					),
				),
			),

			// Features demonstration
			Article(
				Class("card p-4 mb-4"),
				H2(g.Text("Protected Features")),
				P(
					g.Text("This page demonstrates a protected route in the HAGG boilerplate. "),
					g.Text("Only authenticated users can access this content."),
				),

				Div(
					Class("accordion mt-3"),
					ID("featuresAccordion"),

					Div(
						Class("accordion-item"),
						H2(
							Class("accordion-header"),
							Button(
								Class("accordion-button collapsed"),
								Type("button"),
								g.Attr("data-bs-toggle", "collapse"),
								g.Attr("data-bs-target", "#collapseAuth"),
								g.Text("How does authentication work?"),
							),
						),
						Div(
							ID("collapseAuth"),
							Class("accordion-collapse collapse"),
							g.Attr("data-bs-parent", "#featuresAccordion"),
							Div(
								Class("accordion-body"),
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
					),

					Div(
						Class("accordion-item"),
						H2(
							Class("accordion-header"),
							Button(
								Class("accordion-button collapsed"),
								Type("button"),
								g.Attr("data-bs-toggle", "collapse"),
								g.Attr("data-bs-target", "#collapseBuild"),
								g.Text("What can you build with HAGG?"),
							),
						),
						Div(
							ID("collapseBuild"),
							Class("accordion-collapse collapse"),
							g.Attr("data-bs-parent", "#featuresAccordion"),
							Div(
								Class("accordion-body"),
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
				),
			),

			// Call to action
			Article(
				Class("card p-4 text-center"),
				H3(g.Text("Ready to build?")),
				P(
					g.Text("Check out the source code to see how this dashboard is built. "),
					g.Text("All UI is rendered server-side with gomponents."),
				),
				Div(
					Class("mt-3"),
					A(
						Href("https://github.com/axelrhd/hagg"),
						Class("btn btn-outline-secondary"),
						Target("_blank"),
						g.Text("View Source on GitHub"),
					),
				),
			),
		)

		return ctx.Render(layout.Page(ctx, deps, content))
	}
}
