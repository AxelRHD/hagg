// Toast notification system with surreal.js
// Matches the backend toast package and icons.go

// Icon map (must match toast/icons.go exactly)
const TOAST_ICONS = {
    success: `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="10" cy="10" r="9" stroke="#43a047" stroke-width="2"/>
  <path d="M6 10l2.5 2.5L14 7" stroke="#43a047" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`,

    error: `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="10" cy="10" r="9" stroke="#e53935" stroke-width="2"/>
  <path d="M7 7l6 6M13 7l-6 6" stroke="#e53935" stroke-width="2" stroke-linecap="round"/>
</svg>`,

    warning: `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <path d="M10 2L2 17h16L10 2z" stroke="#fb8c00" stroke-width="2" stroke-linejoin="round"/>
  <path d="M10 8v3M10 14h.01" stroke="#fb8c00" stroke-width="2" stroke-linecap="round"/>
</svg>`,

    info: `<svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="10" cy="10" r="9" stroke="#1095c1" stroke-width="2"/>
  <path d="M10 9v5M10 6h.01" stroke="#1095c1" stroke-width="2" stroke-linecap="round"/>
</svg>`
}

// Create toast containers on page load

// Show a toast notification
// Matches the toast.Toast struct from Go
function showToast({ message, level = 'info', timeout = 5000, position = 'bottom-right' }) {
    // Position styles
    const positionStyles = {
        'bottom-right': 'position: fixed; bottom: 1rem; right: 1rem; z-index: 9999;',
        'top-right': 'position: fixed; top: 1rem; right: 1rem; z-index: 9999;',
        'bottom-left': 'position: fixed; bottom: 1rem; left: 1rem; z-index: 9999;',
        'top-left': 'position: fixed; top: 1rem; left: 1rem; z-index: 9999;',
    }

    const icon = TOAST_ICONS[level] || TOAST_ICONS.info
    const toastHtml = `
        <div class="toast toast-${level}" style="opacity: 0; transition: opacity 0.3s; ${positionStyles[position] || positionStyles['bottom-right']}">
            <div style="display: flex; gap: 0.75rem; align-items: center;">
                <div style="flex-shrink: 0;">
                    ${icon}
                </div>
                <div style="flex: 1;">
                    ${escapeHtml(message)}
                </div>
            </div>
        </div>
    `

    // Append directly to body
    document.body.insertAdjacentHTML('beforeend', toastHtml)
    const toast = document.body.lastElementChild

    // Fade in animation
    setTimeout(() => {
        toast.style.opacity = '1'
    }, 10)

    // Auto-remove with fade out
    if (timeout > 0) {
        setTimeout(() => {
            toast.style.opacity = '0'
            setTimeout(() => toast.remove(), 300)
        }, timeout)
    }
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div')
    div.textContent = text
    return div.innerHTML
}
