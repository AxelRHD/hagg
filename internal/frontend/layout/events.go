package layout

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg-lib/hxevents"
	g "maragu.dev/gomponents"
)

// RenderEvents renders initial events as self-destructing HTML elements for full-page loads.
// For HTMX requests, events are sent via HX-Trigger headers instead.
//
// Uses Surreal.js pattern: Each toast is a <div> with showToast() call + me().remove().
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
	// For HTMX requests, don't render DOM elements (events go in headers)
	if hxevents.IsHtmxRequest(ctx.Req.Header) {
		return nil
	}

	// Convert handler.Event to hxevents.Event to avoid import cycle
	hxEvents := make([]hxevents.Event, len(ctx.Events()))
	for i, e := range ctx.Events() {
		hxEvents[i] = hxevents.Event{Name: e.Name, Payload: e.Payload}
	}

	// Render toasts as self-destructing HTML elements
	// Returns nil if no toasts (clean, no empty elements)
	return hxevents.RenderToasts(ctx.Req, hxEvents)
}
