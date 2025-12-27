package shared

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetBasePath(ctx *gin.Context) string {
	bpRaw, ok := ctx.Get("basePath")
	if !ok {
		return ""
	}

	bp, ok := bpRaw.(string)
	if !ok {
		return ""
	}

	if bp == "/" {
		return ""
	}

	return bp
}

func Lnk(ctx *gin.Context, pth string) string {
	return fmt.Sprintf("%v%v", GetBasePath(ctx), pth)
}
