package layout

import (
	"github.com/axelrhd/hagg-lib/notie"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/shared"
	"github.com/gin-gonic/gin"
	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Page(ctx *gin.Context, deps app.Deps, content ...g.Node) g.Node {
	shared.HandleFlash(ctx)

	cn := g.Group{
		// alpine
		x.Data("{ picoTheme: $persist('')}"),
		x.Effect("picoTheme !== '' && document.documentElement.setAttribute('data-theme', picoTheme)"),
		// x.On("notie-alert.document", "notie.alert({ type: $event.detail.type, text: $event.detail.text, position: $event.detail.position, time: $event.detail.time, stay: $event.detail.stay })"),
		x.On("notie-alert.document", "notie.alert({ ...$event.detail })"),

		Div(
			ID("page"),

			shared.MainContainer(
				Navbar(ctx, deps, false),

				Div(
					ID("content"),
					g.Group(content),
				),
			),
		),

		notie.FromContext(ctx),
	}

	return Skeleton(ctx, cn...)
}
