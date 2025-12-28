package layout

import (
	"github.com/axelrhd/hagg-lib/view"
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
				TitleEl(g.Text("HAGG Stack")),

				Charset("utf-8"),
				Meta(
					Name("viewport"),
					Content("width=device-width, initial-scale=1"),
				),

				Meta(
					Name("color-scheme"),
					Content("light dark"),
				),

				// --- ALPINE JS ---
				view.Script(ctx, "/static/js/alpine_persist.min.js",
					Defer(),
				),
				view.Script(ctx, "/static/js/alpine.min.js",
					Defer(),
				),

				// --- HTMX ---
				view.Script(ctx, "/static/js/htmx.min.js"),

				// --- SURREAL JS ---
				view.Script(ctx, "/static/js/surreal_v1.3.4.js"),

				// --- NOTIE ---
				view.Stylesheet(ctx, "/static/css/notie.min.css"),
				view.Script(ctx, "/static/js/notie.min.js"),
				Script(
					g.Raw("notie.setOptions({ positions: { alert: 'bottom', force: 'bottom', confirm: 'bottom', input: 'bottom' } })"),
				),

				// --- FLEXBOXGRID CSS ---
				view.Stylesheet(ctx, "/static/css/flexboxgrid.min.css"),

				// --- TACHYONS CSS ---
				view.Stylesheet(ctx, "/static/css/tachyons.min.css"),

				// --- PICO CSS ---
				view.Stylesheet(ctx, "/static/css/pico.min.css"),
				view.Stylesheet(ctx, "/static/css/pico.colors.min.css"),

				// --- APP CSS ---
				view.Stylesheet(ctx, "/static/css/app.css"),
			),
			Body(
				grp,
			),
		),
	)
}
