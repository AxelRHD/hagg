// Event processing system for HAGG
// Handles both initial-events (full-page loads) and HTMX events (HX-Trigger headers)
// Uses surreal.js for DOM manipulation

// Process initial events on page load
me(document).on('DOMContentLoaded', () => {
    const el = me('#initial-events')
    if (el) {
        try {
            const events = JSON.parse(el.textContent)
            events.forEach(processEvent)
            el.remove()
        } catch (e) {
            console.error('Failed to parse initial events:', e)
        }
    }
})

// Listen for HTMX events (all three phases)
me(document.body).on('htmx:beforeOnLoad', (evt) => {
    const xhr = evt.detail.xhr

    // Process all HX-Trigger headers
    processHxTrigger(xhr.getResponseHeader('HX-Trigger'))
    processHxTrigger(xhr.getResponseHeader('HX-Trigger-After-Swap'))
    processHxTrigger(xhr.getResponseHeader('HX-Trigger-After-Settle'))
})

// Parse and process HX-Trigger header value
function processHxTrigger(headerValue) {
    if (!headerValue) return

    try {
        const events = JSON.parse(headerValue)
        Object.keys(events).forEach(name => {
            processEvent({
                name: name,
                payload: events[name]
            })
        })
    } catch (e) {
        console.error('Failed to parse HX-Trigger:', e)
    }
}

// Central event processor
// Add new event handlers here
function processEvent(event) {
    switch (event.name) {
        case 'toast':
            showToast(event.payload)
            break

        case 'auth-changed':
            // Refresh navigation when auth state changes
            const nav = me('#nav')
            if (nav) {
                htmx.trigger(nav, 'refresh')
            }
            break

        default:
            console.log('Unhandled event:', event.name, event.payload)
            break
    }
}
