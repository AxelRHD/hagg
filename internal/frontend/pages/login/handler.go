package login

import (
	"github.com/axelrhd/hagg-lib/handler"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/shared"
)

// HxLogin handles HTMX login requests.
// It validates the UID, attempts login, and returns appropriate toast notifications.
func HxLogin(deps app.Deps) handler.HandlerFunc {
	return func(ctx *handler.Context) error {
		// Parse form data
		if err := ctx.Req.ParseForm(); err != nil {
			ctx.Toast("Invalid form data").Error().Notify()
			return ctx.NoContent()
		}

		uid := ctx.Req.FormValue("uid")
		if uid == "" {
			ctx.Toast("UID is required").Error().Notify()
			return ctx.NoContent()
		}

		// Attempt login
		_, err := deps.Auth.Login(ctx.Req, uid)
		if err != nil {
			ctx.Toast(err.Error()).Error().Notify()
			return ctx.NoContent()
		}

		// Success
		ctx.Toast("Login erfolgreich.").Success().Notify()
		ctx.Event("auth-changed", true)
		return ctx.NoContent()
	}
}

// HxLogout handles HTMX logout requests.
// It clears the session, sets a flash message, and redirects to home.
func HxLogout(deps app.Deps) handler.HandlerFunc {
	return func(ctx *handler.Context) error {
		err := deps.Auth.Logout(ctx.Req)
		if err != nil {
			ctx.Toast(err.Error()).Error().Notify()
			return ctx.NoContent()
		}

		// Set flash message for redirect
		shared.SetFlash(ctx, "success", "Logout erfolgreich.")

		// HX-Redirect header (must be set before calling NoContent)
		ctx.Res.Header().Set("HX-Redirect", "/")
		return ctx.NoContent()
	}
}
