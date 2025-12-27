package middleware

import (
	"net/http"

	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/flash"
	"github.com/axelrhd/hagg/internal/frontend/shared"
	"github.com/gin-gonic/gin"
)

func RequireAuth(deps app.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if deps.Auth.IsAuthenticated(ctx) {
			ctx.Next()
			return
		}

		flash.Set(ctx, flash.Unauthorized)

		ctx.Redirect(http.StatusFound, shared.Lnk(ctx, "/"))
		ctx.Abort()
	}
}
