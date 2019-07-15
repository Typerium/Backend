package httpserver

import (
	"github.com/gramework/gramework"
)

func notFoundHandler(ctx *gramework.Context) {
	ctx.JSON(NotFoundErr)
}

func methodNotAllowed(ctx *gramework.Context) {
	ctx.JSON(MethodNotAllowedErr)
}
