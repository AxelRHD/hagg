package render

import (
	"log"
	"net/http"

	"github.com/axelrhd/hagg/internal/app"
	"github.com/gin-gonic/gin"
	g "maragu.dev/gomponents"
)

type PageFunc func(*gin.Context, app.Deps) g.Node

func Page(deps app.Deps, fn PageFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := fn(ctx, deps).Render(ctx.Writer); err != nil {
			log.Println(err)
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
