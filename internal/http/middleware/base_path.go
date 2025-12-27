package middleware

import "github.com/gin-gonic/gin"

func BasePath(base string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("basePath", base)
		ctx.Next()
	}
}
