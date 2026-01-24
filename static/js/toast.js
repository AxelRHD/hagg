// Toast notification system using Bootstrap 5.3 Toasts
// Matches the backend toast package interface

// Ensure toast container exists
function getToastContainer() {
    var container = document.getElementById('toastContainer');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toastContainer';
        container.className = 'toast-container position-fixed bottom-0 end-0 p-3';
        document.body.appendChild(container);
    }
    return container;
}

// Show a toast notification
// Matches the toast.Toast struct from Go
function showToast(opts) {
    var message = opts.message || '';
    var level = opts.level || 'info';
    var timeout = opts.timeout !== undefined ? opts.timeout : 5000;

    var container = getToastContainer();
    if (typeof bootstrap === 'undefined') {
        console.error('Bootstrap JS not loaded');
        return;
    }

    var id = 'toast-' + Date.now();

    // Map level to Bootstrap border color class
    var levelMap = {
        success: { border: 'border-success', icon: 'bi-check-circle-fill', color: 'text-success' },
        error:   { border: 'border-danger',  icon: 'bi-x-circle-fill',     color: 'text-danger' },
        warning: { border: 'border-warning', icon: 'bi-exclamation-triangle-fill', color: 'text-warning' },
        info:    { border: 'border-info',    icon: 'bi-info-circle-fill',  color: 'text-info' }
    };

    var config = levelMap[level] || levelMap.info;

    var html = '<div id="' + id + '" class="toast align-items-center bg-body-secondary shadow border-0 border-start border-4 ' + config.border + '" role="alert">' +
        '<div class="d-flex">' +
        '<div class="toast-body">' +
        '<i class="bi ' + config.icon + ' me-2 ' + config.color + '"></i>' +
        escapeHtml(message) +
        '</div>' +
        '<button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast"></button>' +
        '</div></div>';

    container.insertAdjacentHTML('beforeend', html);
    var toastEl = document.getElementById(id);

    var toastOpts = {};
    if (timeout > 0) {
        toastOpts.delay = timeout;
    } else {
        toastOpts.autohide = false;
    }

    var toast = new bootstrap.Toast(toastEl, toastOpts);
    toast.show();

    // Remove from DOM after hidden
    toastEl.addEventListener('hidden.bs.toast', function() {
        toastEl.remove();
    });
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    var div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
