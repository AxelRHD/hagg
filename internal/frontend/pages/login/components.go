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
//	loginURL := view.URLStringChi(req, "/htmx/login")  // Chi
//	LoginForm(loginURL)
func LoginForm(loginURL string) g.Node {
	return Article(
		Class("container-narrow"),

		H1(
			Class("text-center"),
			g.Text("Login"),
		),

		Form(
			hx.Post(loginURL),

			// Listen for toast event dispatched by HTMX from HX-Trigger header
			hx.On("toast", "showToast(event.detail)"),

			Input(
				Type("text"),
				ID("uid"),
				Name("uid"),
				Placeholder("UID (z. B. knO09tSEzCYDhjTcQ…)"),
				Required(),
				AutoFocus(),
			),

			Button(
				Type("submit"),
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
//	logoutURL := view.URLStringChi(req, "/htmx/logout")  // Chi
//	LogoutButton(logoutURL, user.FullName())
func LogoutButton(logoutURL, username string) g.Node {
	return Article(
		Class("container-narrow text-center"),
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
				Class("outline"),
				g.Text("Abmelden"),
			),
		),
	)
}
