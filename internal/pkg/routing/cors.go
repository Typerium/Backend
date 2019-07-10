package routing

import (
	"bytes"
	"strings"

	"github.com/valyala/fasthttp"
)

var accessControlAllowHeaders = strings.Join([]string{
	acceptHeader,
	contentTypeHeader,
	"Content-Length",
	"Accept-Encoding",
	authorizationHeader,
	"X-CSRF-Token",
	requestHeader,
	"X-Apollo-Tracing",
}, ", ")

func (r *Router) corsHandler(ctx *fasthttp.RequestCtx) {
	origin := ctx.Request.Header.Peek(originHeader)
	if origin == nil {
		return
	}
	allowedMethods := "HEAD, OPTIONS"
	route, ok := r.routes[b2s(ctx.Path())]
	if ok {
		allowedMethods = b2s(bytes.Join(route.methods, []byte(", "))) + ", " + allowedMethods
	}
	ctx.Response.Header.SetBytesV(accessControlAllowOriginHeader, origin)
	ctx.Response.Header.Set(accessControlAllowMethods, allowedMethods)
	ctx.Response.Header.Set(accessControlAllowHeadersHeader, accessControlAllowHeaders)
}
