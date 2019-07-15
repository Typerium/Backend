package handlers

import (
	"context"

	"github.com/gramework/gramework"

	"typerium/internal/app/gateway/graphql"
)

func Handler(executor graphql.Executor) gramework.RequestHandler {
	return func(ctx *gramework.Context) {
		req, err := ctx.DecodeGQL()
		if err != nil {
			if err == gramework.ErrInvalidGQLRequest {
				ctx.BadRequest(err)
				return
			}
			ctx.Err500()
			return
		}

		gqlCtx := context.Background()

		resp := executor.Exec(gqlCtx, &graphql.Request{
			Query:         req.Query,
			OperationName: req.OperationName,
			Variables:     req.Variables,
		})

		err = ctx.JSON(resp)
		if err != nil {
			ctx.Err500()
		}
	}
}
