package handlers

import (
	"context"
	"encoding/json"

	"github.com/valyala/fasthttp"

	"typerium/internal/app/gateway/graphql"
	"typerium/internal/pkg/routing"
)

type graphqlHandler struct {
	gqlExecutor graphql.Executor
}

type gqlRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func PostHandler(executor graphql.Executor) routing.Handler {
	return func(ctx routing.Context) {
		ctx.Response().Header().SetContentType(routing.JSONContentType)
		ctx.Response().SetStatusCode(fasthttp.StatusOK)

		var req gqlRequest
		if err := json.Unmarshal(ctx.Request().Body(), &req); err != nil {
			ctx.Response().Error(err)
			return
		}

		graphqlContext := context.Background()

		authHeader := ctx.Request().Authorization()
		if authHeader != nil {
			// 	todo add token to context

		}

		gqlReq := &graphql.Request{
			Query:         req.Query,
			Variables:     req.Variables,
			OperationName: req.OperationName,
		}

		response := executor.Exec(graphqlContext, gqlReq)

		ctx.Response().SetBody(response)
	}
}
