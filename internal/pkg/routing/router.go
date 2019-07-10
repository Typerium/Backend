package routing

import (
	"bytes"
	"fmt"
	"strings"
	"unsafe"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Handler func(ctx Context)

type route struct {
	methods [][]byte
	Handler
}

// http methods
var (
	GET  = []byte("GET")
	POST = []byte("POST")
)

// Router helper tool for request executing
type Router struct {
	routes map[string]*route
	log    *zap.Logger
}

// New constructor for Router
func New(log *zap.Logger) *Router {
	return &Router{
		routes: make(map[string]*route),
		log:    log,
	}
}

// Register add new route to routes map
func (r *Router) Register(path string, handler Handler, methods ...[]byte) {
	if r.routes == nil {
		r.routes = make(map[string]*route)
	}
	r.routes[path] = &route{
		Handler: handler,
		methods: methods,
	}
}

type t struct {
}

// ServeHTTP entry point for fasthttp server
func (r *Router) ServeHTTP(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(xContentTypeOptionsHeader, xContentTypeOptionsNosniff)
	requestID := ctx.Request.Header.Peek(requestHeader)
	if requestID == nil {
		requestID = []byte(fmt.Sprintf("%d", ctx.ID()))
	}

	ctx.Response.Header.SetBytesV(requestHeader, requestID)
	headerFields := make([]zapcore.Field, 0, ctx.Request.Header.Len())
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headerFields = append(headerFields, zap.ByteString(b2s(key), value))
	})
	log := r.log.With(
		zap.ByteString("request_id", requestID),
	).With(headerFields...)

	log.Info(fmt.Sprintf("%s %#v",
		b2s(ctx.URI().FullURI()),
		b2s(ctx.Request.Body())),
	)

	acceptContentType := ctx.Request.Header.Peek(acceptHeader)
	if acceptContentType != nil {
		switch strings.ToLower(b2s(acceptContentType)) {
		case strings.ToLower(JSONContentType):
			ctx.Response.Header.SetContentType(JSONContentType)
		default:
			ctx.Response.Header.SetContentType(TextContentType)
		}
	}

	defer func() {
		r.corsHandler(ctx)
	}()

	if ctx.IsOptions() || ctx.IsHead() {
		return
	}

	route, ok := r.routes[b2s(ctx.Path())]
	if !ok {
		errorHandler(&ctx.Response, notFoundError)
		return
	}

	var isAllowedMethod bool
	for _, method := range route.methods {
		if bytes.Equal(ctx.Method(), method) {
			isAllowedMethod = true
			break
		}
	}
	if !isAllowedMethod {
		errorHandler(&ctx.Response, methodNotAllowedError)
		return
	}

	rCtx := acquireRouterContext(ctx)
	defer releaseRouterContext(rCtx)

	route.Handler(rCtx)
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func b2sPtr(b []byte) *string {
	if b == nil || len(b) == 0 {
		return nil
	}
	return (*string)(unsafe.Pointer(&b))
}
