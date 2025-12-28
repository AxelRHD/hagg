package middleware

import (
	"net/http"

	"github.com/axelrhd/hagg-lib/flash"
	"github.com/axelrhd/hagg-lib/notie"
	"github.com/axelrhd/hagg-lib/view"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/gin-gonic/gin"
)

func RequirePermission(deps app.Deps, action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// ------------------------------------------------------------
		// Authentication fehlt → Flash + Redirect (wie RequireAuth)
		// ------------------------------------------------------------
		user, ok := deps.Auth.CurrentUser(ctx)
		if !ok {
			flash.Set(ctx, flash.Unauthorized)

			ctx.Redirect(http.StatusFound, view.URLString(ctx, "/"))
			ctx.Abort()
			return
		}

		// ------------------------------------------------------------
		// Authorization (Casbin)
		// ------------------------------------------------------------
		allowed, err := deps.Enforcer.Enforce(user.DisplayName, action)
		if err != nil {
			notie.NewAlert("Berechtigungsprüfung fehlgeschlagen (Config Fehler).").
				Error().
				Notify(ctx)

			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !allowed {
			notie.NewAlert("Keine Berechtigung.").
				Warning().
				Notify(ctx)

			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// ------------------------------------------------------------
		// Alles ok → Handler darf ausliefern
		// ------------------------------------------------------------
		ctx.Next()
	}
}
