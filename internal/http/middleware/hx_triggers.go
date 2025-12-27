package middleware

import (
	"encoding/json"

	"github.com/axelrhd/hagg/internal/hxevents"
	"github.com/gin-gonic/gin"
)

func HXTriggers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		commit(ctx, hxevents.Immediate)
		commit(ctx, hxevents.AfterSwap)
		commit(ctx, hxevents.AfterSettle)
	}
}

func commit(ctx *gin.Context, phase hxevents.Phase) {
	raw, ok := ctx.Get(hxevents.CtxKey(phase))
	if !ok {
		return
	}

	b, err := json.Marshal(raw)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header(string(phase), string(b))
}
