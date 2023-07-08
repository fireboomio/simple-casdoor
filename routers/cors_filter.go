package routers

import (
	"net/http"

	"github.com/beego/beego/context"
)

const (
	headerAllowOrigin  = "Access-Control-Allow-Origin"
	headerAllowMethods = "Access-Control-Allow-Methods"
)

func CorsFilter(ctx *context.Context) {
	if ctx.Input.Method() == "OPTIONS" {
		ctx.Output.Header(headerAllowOrigin, "*")
		ctx.Output.Header(headerAllowMethods, "POST, GET, OPTIONS, DELETE")
		ctx.ResponseWriter.WriteHeader(http.StatusOK)
		return
	}
}
