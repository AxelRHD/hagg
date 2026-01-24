package layout

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/shared"
	appshared "github.com/axelrhd/hagg/internal/shared"
	"github.com/axelrhd/hagg/internal/version"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Page renders the full page layout with navbar, content, and event handling.
func Page(ctx *handler.Context, deps app.Deps, content ...g.Node) g.Node {
	// Convert flash messages to toast events
	flashMessages := appshared.GetFlashMessages(ctx)
	for _, msg := range flashMessages {
		switch msg.Level {
		case "success":
			ctx.Toast(msg.Message).Success().Notify()
		case "error":
			ctx.Toast(msg.Message).Error().Notify()
		case "warning":
			ctx.Toast(msg.Message).Warning().Notify()
		case "info":
			ctx.Toast(msg.Message).Info().Notify()
		}
	}

	cn := g.Group{
		Div(
			ID("page"),

			shared.MainContainer(
				Navbar(ctx, deps, false),

				Div(
					ID("content"),
					g.Group(content),
				),

				// Footer with version
				Footer(
					Class("text-center text-body-secondary py-4 mt-5"),
					Small(
						g.Text("HAGG Stack · "),
						Code(Class("text-body-secondary"), g.Text(version.Version)),
					),
				),
			),
		),

		// Render initial events for full-page loads (uses existing RenderEvents from events.go)
		RenderEvents(ctx),
	}

	return Skeleton(ctx.Req, cn...)
}
