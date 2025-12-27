package shared

import (
	"github.com/axelrhd/hagg/internal/flash"
	"github.com/axelrhd/hagg/internal/notie"
	"github.com/gin-gonic/gin"
)

func HandleFlash(ctx *gin.Context) {
	switch {
	case flash.Has(ctx, flash.Unauthorized):
		notie.NewAlert("Bitte einloggen.").
			Warning().
			Notify(ctx)
	case flash.Has(ctx, flash.LogoutSuccess):
		notie.NewAlert("Logout erfolgreich.").
			Success().
			Notify(ctx)
	}
}
