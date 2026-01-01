package login

import (
	"log"
	"net/http"

	"github.com/axelrhd/hagg-lib/flash"
	"github.com/axelrhd/hagg-lib/notie"
	"github.com/axelrhd/hagg-lib/view"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/gin-gonic/gin"
)

func HxLogin(deps app.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		reqData := struct {
			UID string `form:"uid" json:"uid"`
		}{}

		err := ctx.Bind(&reqData)
		if err != nil {
			notie.NewAlert(err.Error()).Error().Notify(ctx)
			log.Println(err)

			ctx.Status(http.StatusNoContent)
			return
		}

		_, err = deps.Auth.Login(ctx, reqData.UID)
		if err != nil {
			notie.NewAlert(err.Error()).Error().Notify(ctx)

			ctx.Status(http.StatusNoContent)
			return
		}

		// hxevents.Add(ctx, hxevents.Immediate, "update-nav", true)
		notie.NewAlert("Login erfolgreich.").Success().Notify(ctx)
		// TODO: Re-enable after migration to Chi/handler.Context
		// hxevents.Add(ctx, hxevents.Immediate, "auth-changed", true)

		ctx.Status(http.StatusNoContent)
	}
}

func HxLogout(deps app.Deps) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := deps.Auth.Logout(ctx)
		if err != nil {
			notie.NewAlert(err.Error()).Error().Notify(ctx)

			ctx.Status(http.StatusNoContent)
			return
		}

		flash.Set(ctx, flash.LogoutSuccess)

		// ctx.Redirect(http.StatusFound, shared.Lnk(ctx, "/"))
		ctx.Header("HX-Redirect", view.URLString(ctx, "/"))
		// ctx.Abort()
	}
}
