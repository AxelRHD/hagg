package layout

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg-lib/hxevents"
	g "maragu.dev/gomponents"
)

// RenderEvents renders initial events script for full-page loads.
// For HTMX requests, events are sent via HX-Trigger headers instead.
//
// This should be called in the layout's Body, before the main content.
//
// Usage:
//
//	html.Body(
//	    RenderEvents(ctx),  // <-- Add this
//	    // ... rest of body content
//	)
func RenderEvents(ctx *handler.Context) g.Node {
	// Convert handler.Event to hxevents.Event to avoid import cycle
	hxEvents := make([]hxevents.Event, len(ctx.Events()))
	for i, e := range ctx.Events() {
		hxEvents[i] = hxevents.Event{Name: e.Name, Payload: e.Payload}
	}

	// Render initial-events script tag
	return hxevents.RenderInitialEvents(ctx.Req, hxEvents)
}
