package middleware

import (
	"net/http"

	"github.com/axelrhd/hagg-lib/casbinx"
	"github.com/axelrhd/hagg-lib/view"
	"github.com/axelrhd/hagg/internal/auth"
	"github.com/axelrhd/hagg/internal/session"
	"github.com/axelrhd/hagg/internal/user"
)

// RequirePermission is a **reference implementation** for Casbin-based authorization.
//
// This middleware combines authentication and authorization in a single check:
//  1. Authentication: Verifies user is logged in (UID in session)
//  2. Authorization: Verifies user has permission via Casbin (subject → action)
//
// # When to use this
//
// Use RequirePermission when your application needs:
//   - Multi-user apps with different roles (admin, editor, viewer)
//   - Fine-grained access control (user:create, user:delete, report:view)
//   - Resource-level permissions beyond "logged in or not"
//
// # When NOT to use this
//
// Skip this middleware and use simpler alternatives when:
//   - Single-user tools: Just use RequireAuth
//   - Simple CRUD apps: Resource ownership checks in handlers are often enough
//   - "Logged in = full access": RequireAuth is sufficient
//
// # Customization
//
// This implementation uses DisplayName as the Casbin subject. Your app may need:
//   - UID as subject (unique, stable)
//   - Email as subject (readable)
//   - Role name directly (if stored in session)
//
// Feel free to adapt or replace this middleware based on your needs.
//
// # Example
//
//	r.Group(func(r chi.Router) {
//	    r.Use(middleware.RequirePermission(deps.Auth, deps.Users, deps.Perms, "dashboard:view"))
//	    r.Get("/dashboard", wrapper.Wrap(dashboard.Page(deps)))
//	})
func RequirePermission(authService *auth.Auth, users user.Store, perms *casbinx.Perm, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCtx := r.Context()

			// Step 1: Check authentication
			rawUID := session.Manager.Get(sessionCtx, auth.SessionKeyUID)
			uid, ok := rawUID.(string)

			if !ok || uid == "" {
				// Not authenticated - redirect to login
				loginURL := view.URLString(r, "/login")

				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", loginURL)
					w.WriteHeader(http.StatusNoContent)
					return
				}

				http.Redirect(w, r, loginURL, http.StatusSeeOther)
				return
			}

			// Step 2: Load user for authorization check
			u, err := users.FindByUID(sessionCtx, uid)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Step 3: Check authorization via Casbin
			// Subject is the user's display name (adapt this if your policy uses UID or email)
			allowed := perms.Can(u.DisplayName, action)
			if !allowed {
				// Not authorized - return 403 with toast for HTMX requests
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Trigger", `{"toast":{"message":"Permission denied.","level":"warning","timeout":3000,"position":"bottom-right"}}`)
					w.WriteHeader(http.StatusNoContent)
					return
				}

				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}

			// User is authenticated and authorized - continue
			next.ServeHTTP(w, r)
		})
	}
}
