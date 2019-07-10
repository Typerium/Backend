package routing

import (
	"sync"

	"github.com/valyala/fasthttp"
)

// Context custom context for handler
type Context interface {
	Request() Request
	Response() Response
}

// Header access to http header
type Header interface {
	Peek(key string) []byte
	SetContentType(contentType string)
}

// Request actions on input request
type Request interface {
	Header() Header
	Body() []byte
	Authorization() *string
}

// Response actions on output response
type Response interface {
	Header() Header
	SetStatusCode(code int)
	SetBody(body []byte)
	Error(err error)
}

var (
	routerContextPool   sync.Pool
	requestContextPool  sync.Pool
	responseContextPool sync.Pool
)

func init() {
	routerContextPool = sync.Pool{
		New: func() interface{} {
			return new(routerContext)
		},
	}
	requestContextPool = sync.Pool{
		New: func() interface{} {
			return new(request)
		},
	}
	responseContextPool = sync.Pool{
		New: func() interface{} {
			return new(response)
		},
	}
}

func acquireRouterContext(ctx *fasthttp.RequestCtx) *routerContext {
	v, ok := routerContextPool.Get().(*routerContext)
	if !ok || v == nil {
		v = new(routerContext)
	}

	v.req = acquireRequestContext(&ctx.Request)
	v.resp = acquireResponseContext(&ctx.Response)
	return v
}

func releaseRouterContext(ctx *routerContext) {
	releaseRequestContext(ctx.req)
	releaseResponseContext(ctx.resp)
	routerContextPool.Put(ctx)
}

func acquireRequestContext(req *fasthttp.Request) *request {
	v, ok := requestContextPool.Get().(*request)
	if !ok || v == nil {
		v = new(request)
	}
	v.Request = req
	return v
}

func releaseRequestContext(v *request) {
	requestContextPool.Put(v)
}

func acquireResponseContext(resp *fasthttp.Response) *response {
	v, ok := responseContextPool.Get().(*response)
	if !ok || v == nil {
		v = new(response)
	}
	v.Response = resp
	return v
}

func releaseResponseContext(v *response) {
	responseContextPool.Put(v)
}

type routerContext struct {
	req  *request
	resp *response
}

type request struct {
	*fasthttp.Request
}

func (req *request) Authorization() *string {
	return b2sPtr(req.Request.Header.Peek(authorizationHeader))
}

func (req *request) Header() Header {
	return &req.Request.Header
}

type response struct {
	*fasthttp.Response
}

func (resp *response) Header() Header {
	return &resp.Response.Header
}

func (resp *response) Error(err error) {
	httpErr, ok := err.(*httpError)
	if !ok {
		httpErr = internalError
	}
	errorHandler(resp.Response, httpErr)
}

func (ctx *routerContext) Request() Request {
	return ctx.req
}

func (ctx *routerContext) Response() Response {
	return ctx.resp
}
