package login

import (
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

// LoginForm renders a simple UID-only login form.
// Framework-agnostic - accepts URL string instead of context.
//
// Usage:
//
//	loginURL := view.URLString(ctx, "/htmx/login")  // Gin
//	loginURL := view.URLString(req, "/htmx/login")  // Chi
//	LoginForm(loginURL)
func LoginForm(loginURL string) g.Node {
	return Article(
		Class("container-narrow card p-4"),

		H1(
			Class("text-center mb-4"),
			g.Text("Login"),
		),

		Form(
			hx.Post(loginURL),

			Div(
				Class("mb-3"),
				Input(
					Type("password"),
					Class("form-control"),
					ID("uid"),
					Name("uid"),
					Placeholder("UID"),
					Required(),
					AutoFocus(),
				),
			),

			Button(
				Type("submit"),
				Class("btn btn-primary w-100"),
				g.Text("Anmelden"),
			),
		),
	)
}

// LogoutButton renders a logout button with username display.
// Framework-agnostic - accepts URL string and username.
//
// Usage:
//
//	logoutURL := view.URLString(ctx, "/htmx/logout")  // Gin
//	logoutURL := view.URLString(req, "/htmx/logout")  // Chi
//	LogoutButton(logoutURL, user.FullName())
func LogoutButton(logoutURL, username string) g.Node {
	return Article(
		Class("container-narrow text-center card p-4"),
		H3(g.Text("Bereits eingeloggt")),
		P(
			g.Text("Du bist angemeldet als "),
			Strong(g.Text(username)),
			g.Text("."),
		),
		Form(
			hx.Post(logoutURL),
			Style("margin-top: 1.5rem;"),
			Button(
				Type("submit"),
				Class("btn btn-outline-secondary"),
				g.Text("Abmelden"),
			),
		),
	)
}
