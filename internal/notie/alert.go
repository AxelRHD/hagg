package notie

import (
	"fmt"

	"github.com/axelrhd/hagg/internal/hxevents"
	"github.com/gin-gonic/gin"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const (
	HXEventNotieAlert = "notie-alert"
	CtxNotieAlert     = "notie-alert"
)

type Alert struct {
	Text     string `json:"text"`
	Level    string `json:"type"`     // default: info ['success', 'warning', 'error', 'info', 'neutral']
	Stay     bool   `json:"stay"`     // default: false
	Time     int    `json:"time"`     // default: 3, min: 1
	Position string `json:"position"` // default: top ['top', 'bottom']
}

func NewAlert(msg string) *Alert {
	return &Alert{
		Text:     msg,
		Level:    "info",
		Stay:     false,
		Time:     3,
		Position: "bottom",
	}
}

func (a *Alert) Success() *Alert {
	a.Level = "success"

	return a
}

func (a *Alert) Warning() *Alert {
	a.Level = "warning"

	return a
}

func (a *Alert) Error() *Alert {
	a.Level = "error"

	return a
}

func (a *Alert) Info() *Alert {
	a.Level = "info"

	return a
}

func (a *Alert) Neutral() *Alert {
	a.Level = "neutral"

	return a
}

func (a *Alert) SetTimeout(dur int) *Alert {
	a.Time = dur

	return a
}

func (a *Alert) ShowTop() *Alert {
	a.Position = "top"

	return a
}

func (a *Alert) ShowBottom() *Alert {
	a.Position = "bottom"

	return a
}

func (a Alert) Gomp() g.Node {
	return Div(
		Class("notie-alert"),
		Script(
			g.Raw(a.Script()),
		),
		Script(
			g.Raw("me().remove()"),
		),
	)
}

func (a Alert) String() string {
	return fmt.Sprint(a.Gomp())
}

func (a Alert) Script() string {
	return fmt.Sprintf("notie.alert({ text: '%s', type: '%s', stay: %v, time: %v, position: '%s' })",
		a.Text, a.Level, a.Stay, a.Time, a.Position)

}

func (a *Alert) HXTrigger(ctx *gin.Context) {
	hxevents.Add(ctx, hxevents.Immediate, HXEventNotieAlert, a.Payload())
}

func (a *Alert) Payload() map[string]any {
	return map[string]any{
		"text": a.Text,
		"type": a.Level,
		"stay": a.Stay,
		"time": a.Time,
	}
}

func (a *Alert) Notify(ctx *gin.Context) {
	if ctx.GetHeader("HX-Request") == "true" {
		a.HXTrigger(ctx)
		return
	}

	ctx.Set(CtxNotieAlert, a)
}

func FromContext(ctx *gin.Context) g.Node {
	raw, ok := ctx.Get(CtxNotieAlert)
	if !ok {
		return nil
	}

	alert, ok := raw.(*Alert)
	if !ok {
		return nil
	}

	return alert.Gomp()
}
