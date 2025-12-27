package login

import (
	"log"
	"net/http"

	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/flash"
	"github.com/axelrhd/hagg/internal/frontend/shared"
	"github.com/axelrhd/hagg/internal/hxevents"
	"github.com/axelrhd/hagg/internal/notie"
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
		hxevents.Add(ctx, hxevents.Immediate, "auth-changed", true)

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
		ctx.Header("HX-Redirect", shared.Lnk(ctx, "/"))
		// ctx.Abort()
	}
}
