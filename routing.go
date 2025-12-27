package hagg

import (
	"github.com/gin-gonic/gin"
	"github.com/axelrhd/hagg/internal/app"
	"github.com/axelrhd/hagg/internal/frontend/pages/login"
	"github.com/axelrhd/hagg/internal/http/render"
)

func Routing(rg *gin.RouterGroup, deps app.Deps) {
	mGetPost := []string{"GET", "POST"}

	// pages
	rg.Match(mGetPost, "/", render.Page(deps, login.Page))

	// htmx-endpoints
	hxg := rg.Group("/htmx")

	hxg.POST("/login", login.HxLogin(deps))
	hxg.POST("/logout", login.HxLogout(deps))
}
