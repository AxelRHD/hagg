package login

import (
	"github.com/axelrhd/hagg-lib/view"
	"github.com/gin-gonic/gin"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

// LoginForm renders a simple UID-only login form.
// Expects POST /login with field "uid".
func LoginForm(ctx *gin.Context) g.Node {
	return Article(
		// Pico card + Tachyons width
		Class("w-100 mw6"),

		H1(
			Class("tc"),
			g.Text("Login"),
		),

		Form(
			hx.Post(view.URLString(ctx, "/htmx/login")),

			Input(
				Type("text"),
				ID("uid"),
				Name("uid"),
				Placeholder("UID (z. B. knO09tSEzCYDhjTcQ…)"),
				Required(),
				AutoFocus(),
			),
			// Submit button
			Div(
				Class("mt3"),

				Button(
					Type("submit"),
					Class("w-100"),
					g.Text("Anmelden"),
				),
			),
		),
	)
}
