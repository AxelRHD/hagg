package middleware

import (
	"net/http"

	"github.com/axelrhd/hagg-lib/handler"
)

// Logger is a Chi-compatible middleware that logs HTTP requests.
// It logs the request method and path using the wrapper's logger.
//
// Example:
//
//	r.Use(middleware.Logger(wrapper))
func Logger(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapper.Logger().Info("request",
				"method", r.Method,
				"path", r.URL.Path,
			)
			next.ServeHTTP(w, r)
		})
	}
}

// Recovery is a Chi-compatible middleware that recovers from panics.
// It logs the panic and returns a 500 Internal Server Error.
//
// Example:
//
//	r.Use(middleware.Recovery(wrapper))
func Recovery(wrapper *handler.Wrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					wrapper.Logger().Error("panic recovered",
						"error", err,
						"path", r.URL.Path,
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// CORS is a Chi-compatible middleware that sets CORS headers.
// It allows all origins and common HTTP methods + HTMX headers.
//
// Example:
//
//	r.Use(middleware.CORS())
func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, HX-Request, HX-Trigger, HX-Target")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimit is a placeholder middleware for rate limiting.
// TODO: Implement actual rate limiting logic
//
// Example:
//
//	r.Use(middleware.RateLimit)
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement rate limiting
		next.ServeHTTP(w, r)
	})
}

// Secure is a Chi-compatible middleware that sets security headers.
// It sets common security headers to protect against XSS, clickjacking, etc.
//
// Example:
//
//	r.Use(middleware.Secure)
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
