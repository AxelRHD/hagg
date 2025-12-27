package login

import (
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/layout"
	"github.com/axelrhd/hagg/internal/frontend/shared"
	"log"

	"github.com/gin-gonic/gin"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func Page(ctx *gin.Context, deps app.Deps) g.Node {
	sub := "arudolf" // oder ctx.MustGet(...) → aktuell fest verdrahtet
	ok, err := deps.Enforcer.Enforce(sub, "user:list")
	if err != nil {
		log.Println("casbin error:", err)
	} else {
		log.Println("casbin user:list allowed =", ok)
	}

	uid := ""
	user, _ := deps.Auth.CurrentUser(ctx)
	if user != nil {
		uid = user.UID
	}

	content := g.Group{
		Div(
			hx.Post(shared.Lnk(ctx, "/")),
			hx.Trigger("auth-changed from:body"),
			hx.Target("#page"),
			hx.Select("#page"),
			hx.Swap("outerHTML"),

			// Centering
			Class("flex items-center justify-center pa3"),
			Style("min-height: 80vh"),

			g.If(!deps.Auth.IsAuthenticated(ctx),
				LoginForm(ctx),
			),

			g.If(deps.Auth.IsAuthenticated(ctx),
				Article(
					Class("w-100 mw6 tc"),

					H3(g.Text("Bereits eingeloggt")),

					P(
						Class("f6 mt2"),
						g.Text("Du bist angemeldet als "),
						Strong(g.Text(uid)),
						g.Text("."),
					),

					Form(
						hx.Post(shared.Lnk(ctx, "/htmx/logout")),
						Class("mt3"),

						Button(
							Type("submit"),
							Class("outline"),
							g.Text("Abmelden"),
						),
					),
				),
			),
		),
	}

	return layout.Page(ctx, deps, content...)
}
