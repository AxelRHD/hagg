package layout

import (
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/shared"
	lucide "github.com/eduardolat/gomponents-lucide"
	"github.com/gin-gonic/gin"
	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func Navbar(ctx *gin.Context, deps app.Deps, oob bool) g.Node {
	return Nav(
		ID("navbar"),
		g.If(oob,
			hx.SwapOOB("navbar"),
		),

		Ul(
			Li(
				Strong(
					// Class("nav-title"),
					Class("courier ttu"),
					Style("font-size: 2rem; letter-spacing: 0em"),
					g.Text("KL - Werkzeugkasten"),
				),
			),
		),

		Ul(
			Li(
				A(
					Href(shared.Lnk(ctx, "/")),
					g.Text("Home"),
				),
			),
			g.If(deps.Auth.IsAuthenticated(ctx),
				Li(
					A(
						Href("#"),
						g.Text("Products"),
					),
				),
			),
			g.If(deps.Auth.IsAuthenticated(ctx),
				Li(
					Button(
						// Class("outline"),
						x.Bind("class", "picoTheme === 'dark' ? 'contrast' : 'secondary'"),

						x.On("click.prevent", "console.log('Logout clicked')"),
						hx.Post(shared.Lnk(ctx, "/htmx/logout")),

						g.Text("Logout"),
					),
				),
			),
			Li(
				Button(
					Class("theme-toggle bg-transparent b--none ph0"),

					x.On("click", "picoTheme = picoTheme === 'light' ? 'dark' : 'light'"),

					lucide.Moon(
						x.Show("picoTheme === 'light'"),
					),
					lucide.Sun(
						x.Show("picoTheme === 'dark'"),
					),
				),
			),
		),
	)
}
