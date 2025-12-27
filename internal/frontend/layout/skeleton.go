package layout

import (
	"github.com/axelrhd/hagg/internal/frontend/shared"
	"github.com/gin-gonic/gin"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Skeleton(ctx *gin.Context, content ...g.Node) g.Node {
	grp := g.Group(content)

	return Doctype(
		HTML(
			Lang("en"),
			Data("theme", "dark"),

			Head(
				Charset("utf-8"),
				Meta(
					Name("viewport"),
					Content("width=device-width, initial-scale=1"),
				),

				Meta(
					Name("color-scheme"),
					Content("light dark"),
				),

				TitleEl(g.Text("HAGG Stack")),

				// --- ALPINE JS ---
				Script(
					Defer(),
					Src(shared.Lnk(ctx, "/static/js/alpine_persist.min.js")),
				),
				Script(
					Defer(),
					Src(shared.Lnk(ctx, "/static/js/alpine.min.js")),
				),

				// --- HTMX ---
				Script(
					Src(shared.Lnk(ctx, "/static/js/htmx.min.js")),
				),

				// --- SURREAL JS ---
				Script(
					Src(shared.Lnk(ctx, "/static/js/surreal_v1.3.4.js")),
				),

				// --- NOTIE ---
				Link(
					Rel("stylesheet"),
					Href(shared.Lnk(ctx, "/static/css/notie.min.css")),
				),
				Script(
					Src(shared.Lnk(ctx, "/static/js/notie.min.js")),
				),
				Script(
					g.Raw("notie.setOptions({ positions: { alert: 'bottom', force: 'bottom', confirm: 'bottom', input: 'bottom' } })"),
				),

				// --- FLEXBOXGRID CSS ---
				Link(
					Rel("stylesheet"),
					Href(shared.Lnk(ctx, "/static/css/flexboxgrid.min.css")),
				),

				// --- TACHYONS CSS ---
				Link(
					Rel("stylesheet"),
					Href(shared.Lnk(ctx, "/static/css/tachyons.min.css")),
				),

				// --- PICO CSS ---
				Link(
					Rel("stylesheet"),
					Href(shared.Lnk(ctx, "/static/css/pico.min.css")),
				),
				Link(
					Rel("stylesheet"),
					Href(shared.Lnk(ctx, "/static/css/pico.colors.min.css")),
				),

				// --- APP CSS ---
				Link(
					Rel("stylesheet"),
					Href(shared.Lnk(ctx, "/static/css/app.css")),
				),
			),
			Body(
				grp,
			),
		),
	)
}
