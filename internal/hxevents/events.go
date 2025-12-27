package hxevents

import "github.com/gin-gonic/gin"

type Phase string

const (
	Immediate   Phase = "HX-Trigger"
	AfterSwap   Phase = "HX-Trigger-After-Swap"
	AfterSettle Phase = "HX-Trigger-After-Settle"
)

func ctxKey(p Phase) string {
	return "hx-events:" + string(p)
}

// CtxKey wird von der Middleware verwendet
func CtxKey(p Phase) string {
	return ctxKey(p)
}

// Add registriert ein HX-Event für eine bestimmte Phase
func Add(ctx *gin.Context, phase Phase, name string, payload any) {
	key := ctxKey(phase)

	raw, ok := ctx.Get(key)
	if !ok {
		ctx.Set(key, map[string]any{
			name: payload,
		})
		return
	}

	events := raw.(map[string]any)
	events[name] = payload
}
