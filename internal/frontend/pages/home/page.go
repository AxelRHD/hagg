package home

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/layout"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Page renders the homepage with HAGG stack presentation.
func Page(ctx *handler.Context, deps app.Deps) error {
	content := Div(
		Class("container"),
		Style("padding: 2rem 0;"),

		// Hero section
		Header(
			Class("text-center"),
			Style("margin-bottom: 3rem;"),

			H1(
				Style("font-size: 3rem; margin-bottom: 1rem;"),
				g.Text("HAGG Stack"),
			),
			P(
				Style("font-size: 1.25rem; color: var(--muted-color);"),
				g.Text("A modern, minimal Go web application boilerplate"),
			),
		),

		// What is HAGG section
		Article(
			Style("margin-bottom: 2rem;"),
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
			Style("margin-bottom: 2rem;"),
			H2(g.Text("Features")),
			Div(
				Class("grid"),

				Div(
					H3(g.Text("🎯 Type-Safe HTML")),
					P(g.Text("Write HTML in Go with gomponents. No template files, full compile-time safety.")),
				),

				Div(
					H3(g.Text("⚡ HTMX Hypermedia")),
					P(g.Text("Server-driven UI updates without heavy JavaScript frameworks.")),
				),

				Div(
					H3(g.Text("🔐 Built-in Auth")),
					P(g.Text("Session-based authentication with SCS and SQLite persistence.")),
				),

				Div(
					H3(g.Text("🎨 Dark Mode")),
					P(g.Text("Pico-inspired CSS with Tailwind v4. Dark mode with Alpine.js persistence.")),
				),

				Div(
					H3(g.Text("📦 Event System")),
					P(g.Text("Server-to-client events via HX-Trigger headers. Toast notifications built-in.")),
				),

				Div(
					H3(g.Text("🚀 Minimal Dependencies")),
					P(g.Text("Chi router, SCS sessions, Gomponents. No npm, no build step (except CSS).")),
				),
			),
		),

		// Tech Stack section
		Article(
			Style("margin-bottom: 2rem;"),
			H2(g.Text("Tech Stack")),
			P(g.Text("This boilerplate demonstrates a complete HAGG application with:")),
			Ul(
				Li(Strong(g.Text("Backend:")), g.Text(" Go 1.23+, Chi v5 router, SCS v2 sessions, SQLite")),
				Li(Strong(g.Text("Frontend:")), g.Text(" HTMX, Alpine.js, Surreal.js, Tailwind CSS v4")),
				Li(Strong(g.Text("Templates:")), g.Text(" Gomponents (type-safe HTML in Go)")),
				Li(Strong(g.Text("Authorization:")), g.Text(" Casbin (RBAC/ABAC ready)")),
				Li(Strong(g.Text("Dev Tools:")), g.Text(" Just (task runner), Air (hot reload)")),
			),
		),

		// Get Started section
		Article(
			H2(g.Text("Try it out!")),
			P(
				g.Text("This boilerplate includes authentication and protected routes. "),
				g.Text("Try logging in to see the full features."),
			),
			Div(
				Style("margin-top: 1.5rem; display: flex; gap: 1rem; justify-content: center;"),

				g.If(!deps.Auth.IsAuthenticated(ctx.Req),
					A(
						Href("/login"),
						Role("button"),
						g.Text("Go to Login"),
					),
				),

				g.If(deps.Auth.IsAuthenticated(ctx.Req),
					A(
						Href("/dashboard"),
						Role("button"),
						g.Text("Go to Dashboard"),
					),
				),

				A(
					Href("https://github.com/axelrhd/hagg"),
					Role("button"),
					Class("outline"),
					g.Text("View on GitHub"),
				),
			),
		),
	)

	return ctx.Render(layout.Page(ctx, deps, content))
}
