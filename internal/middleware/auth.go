package middleware

import (
	"net/http"

	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/session"
)

// RequireAuth is a Chi-compatible middleware that requires authentication.
// If the user is not authenticated (user_id not in session), it redirects to the login page.
//
// Example:
//
//	r.Group(func(r chi.Router) {
//	    r.Use(middleware.RequireAuth(wrapper))
//	    r.Get("/dashboard", wrapper.Wrap(dashboardHandler.Index))
//	})
func RequireAuth(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCtx := r.Context()
			userID := session.Manager.GetInt(sessionCtx, "user_id")

			if userID == 0 {
				// Not authenticated - redirect to login
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			// Authenticated - continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RequireGuest is a Chi-compatible middleware that requires the user to NOT be authenticated.
// If the user is already authenticated, it redirects to the home page.
// This is useful for login/register pages.
//
// Example:
//
//	r.Group(func(r chi.Router) {
//	    r.Use(middleware.RequireGuest(wrapper))
//	    r.Get("/auth/login", wrapper.Wrap(authHandler.LoginPage))
//	    r.Post("/auth/login", wrapper.Wrap(authHandler.Login))
//	})
func RequireGuest(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCtx := r.Context()
			userID := session.Manager.GetInt(sessionCtx, "user_id")

			if userID != 0 {
				// Already authenticated - redirect to home
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			// Not authenticated - continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}
