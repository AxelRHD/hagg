package layout

import (
	"net/http"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

// Skeleton renders the HTML document structure.
func Skeleton(req *http.Request, content ...g.Node) g.Node {
	grp := g.Group(content)

	return Doctype(
		HTML(
			Lang("en"),
			// Bootstrap 5.3 dark mode - default dark
			g.Attr("data-bs-theme", "dark"),

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

				// CRITICAL: Inline script to prevent theme flickering
				// This must run BEFORE CSS loads to avoid FOUC (Flash of Unstyled Content)
				// Alpine.js persist key format: _x_theme
				g.Raw(`<script>
					(function() {
						try {
							const stored = localStorage.getItem('_x_theme');
							if (stored) {
								const theme = JSON.parse(stored);
								if (theme) {
									document.documentElement.setAttribute('data-bs-theme', theme);
								}
							}
						} catch (e) {
							// Ignore localStorage errors
						}
					})();
				</script>`),

				// --- BOOTSTRAP 5.3 CSS ---
				Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css")),

				// --- BOOTSTRAP ICONS ---
				Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css")),

				// --- ALPINE JS ---
				Script(
					Src("/static/js/alpine_persist.min.js"),
					Defer(),
				),
				Script(
					Src("/static/js/alpine.min.js"),
					Defer(),
				),

				// --- HTMX ---
				Script(Src("/static/js/htmx.min.js")),

				// --- SURREAL JS ---
				Script(Src("/static/js/surreal_v1.3.4.js")),

				// --- TOAST JS ---
				Script(Src("/static/js/toast.js")),

				// --- APP CSS (custom overrides) ---
				Link(Rel("stylesheet"), Href("/static/css/app.css")),
			),
			Body(
				// Alpine.js state for theme toggle - must be on body tag
				g.Attr("x-data", "{ theme: $persist('dark') }"),
				g.Attr("x-effect", "theme !== '' && document.documentElement.setAttribute('data-bs-theme', theme)"),

				// Global HTMX toast listener - catches toast events from ANY HTMX request
				// IMPORTANT: Must be on <body>, not on individual forms!
				// See CLAUDE.md "Toast System" for details.
				hx.On("toast", "showToast(event.detail)"),

				grp,

				// --- BOOTSTRAP 5.3 JS BUNDLE (includes Popper) ---
				Script(Src("https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js")),
			),
		),
	)
}
