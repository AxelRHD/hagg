package flash

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Key ist ein typisierter Session-Key für Flash-Flags.
type Key string

const (
	Unauthorized  Key = "flash-unauthorized"
	LogoutSuccess Key = "flash-logout-successful"
)

// Set setzt ein Flash-Flag (immer true).
func Set(ctx *gin.Context, k Key) {
	sess := sessions.Default(ctx)
	sess.Set(string(k), true)
	_ = sess.Save()
}

// Has prüft, ob das Flash-Flag existiert.
// Wenn ja, wird es sofort gelöscht (one-shot).
func Has(ctx *gin.Context, k Key) bool {
	sess := sessions.Default(ctx)

	v := sess.Get(string(k))
	ok := v == true

	if ok {
		sess.Delete(string(k))
		_ = sess.Save()
	}

	return ok
}

// Clear löscht das Flash-Flag ohne es zu konsumieren.
func Clear(ctx *gin.Context, k Key) {
	sess := sessions.Default(ctx)
	sess.Delete(string(k))
	_ = sess.Save()
}
