package home

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/layout"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Page renders the homepage with HAGG stack presentation.
func Page(deps app.Deps) handler.HandlerFunc {
	return func(ctx *handler.Context) error {
		content := Div(
			Class("container py-4"),

			// Hero section
			Header(
				Class("text-center mb-5"),

				H1(
					Class("display-4 mb-3"),
					g.Text("HAGG Stack"),
				),
				P(
					Class("lead text-body-secondary"),
					g.Text("A modern, minimal Go web application boilerplate"),
				),
			),

			// What is HAGG section
			Article(
				Class("card mb-4 p-4"),
				H2(g.Text("What is HAGG?")),
				P(
					g.Text("HAGG is a full-stack web development approach combining:"),
				),
				Ul(
					Li(Strong(g.Text("H")), g.Text("TMX – Hypermedia-driven interactions")),
					Li(Strong(g.Text("A")), g.Text("lpine.js – Reactive client-side behavior")),
					Li(Strong(g.Text("G")), g.Text("omponents – Type-safe HTML in Go")),
					Li(Strong(g.Text("G")), g.Text("o – Server-side rendering with Chi router")),
				),
			),

			// Features section
			Article(
				Class("mb-4"),
				H2(Class("mb-4"), g.Text("Features")),
				Div(
					Class("row g-4"),

					Div(
						Class("col-md-6 col-lg-4"),
						Div(
							Class("card h-100 p-3"),
							H3(Class("h5"), g.Text("🎯 Type-Safe HTML")),
							P(Class("mb-0"), g.Text("Write HTML in Go with gomponents. No template files, full compile-time safety.")),
						),
					),

					Div(
						Class("col-md-6 col-lg-4"),
						Div(
							Class("card h-100 p-3"),
							H3(Class("h5"), g.Text("⚡ HTMX Hypermedia")),
							P(Class("mb-0"), g.Text("Server-driven UI updates without heavy JavaScript frameworks.")),
						),
					),

					Div(
						Class("col-md-6 col-lg-4"),
						Div(
							Class("card h-100 p-3"),
							H3(Class("h5"), g.Text("🔐 Built-in Auth")),
							P(Class("mb-0"), g.Text("Session-based authentication with SCS and SQLite persistence.")),
						),
					),

					Div(
						Class("col-md-6 col-lg-4"),
						Div(
							Class("card h-100 p-3"),
							H3(Class("h5"), g.Text("🎨 Dark Mode")),
							P(Class("mb-0"), g.Text("Bootstrap 5.3 with native dark mode support and Alpine.js persistence.")),
						),
					),

					Div(
						Class("col-md-6 col-lg-4"),
						Div(
							Class("card h-100 p-3"),
							H3(Class("h5"), g.Text("📦 Event System")),
							P(Class("mb-0"), g.Text("Server-to-client events via HX-Trigger headers. Toast notifications built-in.")),
						),
					),

					Div(
						Class("col-md-6 col-lg-4"),
						Div(
							Class("card h-100 p-3"),
							H3(Class("h5"), g.Text("🚀 Minimal Dependencies")),
							P(Class("mb-0"), g.Text("Chi router, SCS sessions, Gomponents. No npm, no build step required.")),
						),
					),
				),
			),

			// Tech Stack section
			Article(
				Class("card mb-4 p-4"),
				H2(g.Text("Tech Stack")),
				P(g.Text("This boilerplate demonstrates a complete HAGG application with:")),
				Ul(
					Li(Strong(g.Text("Backend:")), g.Text(" Go 1.23+, Chi v5 router, SCS v2 sessions, SQLite")),
					Li(Strong(g.Text("Frontend:")), g.Text(" HTMX, Alpine.js, Surreal.js, Bootstrap 5.3")),
					Li(Strong(g.Text("Templates:")), g.Text(" Gomponents (type-safe HTML in Go)")),
					Li(Strong(g.Text("Authorization:")), g.Text(" Casbin (RBAC/ABAC ready)")),
					Li(Strong(g.Text("Dev Tools:")), g.Text(" Just (task runner), Air (hot reload)")),
				),
			),

			// Get Started section
			Article(
				Class("card p-4 text-center"),
				H2(g.Text("Try it out!")),
				P(
					g.Text("This boilerplate includes authentication and protected routes. "),
					g.Text("Try logging in to see the full features."),
				),
				Div(
					Class("d-flex gap-2 justify-content-center flex-wrap mt-3"),

					g.If(!deps.Auth.IsAuthenticated(ctx.Req),
						A(
							Href("/login"),
							Class("btn btn-primary"),
							g.Text("Go to Login"),
						),
					),

					g.If(deps.Auth.IsAuthenticated(ctx.Req),
						A(
							Href("/dashboard"),
							Class("btn btn-primary"),
							g.Text("Go to Dashboard"),
						),
					),

					A(
						Href("https://github.com/axelrhd/hagg"),
						Class("btn btn-outline-secondary"),
						g.Text("View on GitHub"),
					),
				),
			),

		)

		return ctx.Render(layout.Page(ctx, deps, content))
	}
}
