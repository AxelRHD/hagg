package middleware

import (
	"net/http"

	"github.com/axelrhd/hagg-lib/flash"
	"github.com/axelrhd/hagg-lib/view"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/gin-gonic/gin"
)

func RequireAuth(deps app.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if deps.Auth.IsAuthenticated(ctx) {
			ctx.Next()
			return
		}

		flash.Set(ctx, flash.Unauthorized)

		ctx.Redirect(http.StatusFound, view.URLString(ctx, "/"))
		ctx.Abort()
	}
}
