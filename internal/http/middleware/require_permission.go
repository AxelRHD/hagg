package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/flash"
	"github.com/axelrhd/hagg/internal/frontend/shared"
	"github.com/axelrhd/hagg/internal/notie"
)

func RequirePermission(deps app.Deps, action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// ------------------------------------------------------------
		// Authentication fehlt → Flash + Redirect (wie RequireAuth)
		// ------------------------------------------------------------
		user, ok := deps.Auth.CurrentUser(ctx)
		if !ok {
			flash.Set(ctx, flash.Unauthorized)

			ctx.Redirect(http.StatusFound, shared.Lnk(ctx, "/"))
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
