package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeployRootRoutes(r gin.IRouter) {
	r.GET("", func(ctx *gin.Context) {
		ctx.Header("Location", "/web/index")
		ctx.Status(http.StatusMovedPermanently)
	})
	r.GET("index", func(ctx *gin.Context) {
		ctx.Header("Location", "/web/index")
		ctx.Status(http.StatusMovedPermanently)
	})
	r.GET("web/null", func(ctx *gin.Context) {
		ctx.File("./static/null")
	})
	r.GET("null", func(ctx *gin.Context) {
		ctx.File("./static/null")
	})
	r.GET("favicon.ico", func(ctx *gin.Context) {
		ctx.File("./static/favicon.ico")
	})
	r.GET("/web/favicon.ico", func(ctx *gin.Context) {
		ctx.File("./static/favicon.ico")
	})
}
