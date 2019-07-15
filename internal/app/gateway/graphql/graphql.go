package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/gqlerror"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"typerium/internal/app/gateway/graphql/internal/response"
	"typerium/internal/app/gateway/graphql/internal/server"
)

//go:generate go run github.com/99designs/gqlgen --config gqlgen.yml

type Request struct {
	Query         string
	OperationName string
	Variables     map[string]interface{}
}

type Executor interface {
	Exec(ctx context.Context, request *Request) *graphql.Response
}

func New(log *zap.Logger) Executor {

	return &gqlGenExecutor{
		ExecutableSchema: server.NewExecutableSchema(server.Config{
			Resolvers: &server.Resolver{},
		}),
		log: log.Named("graphql"),
	}
}

type gqlGenExecutor struct {
	graphql.ExecutableSchema
	log *zap.Logger
}

func (executor *gqlGenExecutor) Exec(ctx context.Context, req *Request) (resp *graphql.Response) {
	defer func() {
		gqlErrors := make([]*gqlerror.Error, 0, len(resp.Errors))
		for index := range resp.Errors {
			if resp.Errors[index] == nil {
				continue
			}
			gqlErrors = append(gqlErrors, resp.Errors[index])
		}
		resp.Errors = gqlErrors
	}()

	tracer := newTracer()

	ctx = tracer.StartOperationExecution(ctx)

	reqCtx := &graphql.RequestContext{
		RawQuery:            req.Query,
		ResolverMiddleware:  graphql.DefaultResolverMiddleware,
		DirectiveMiddleware: graphql.DefaultDirectiveMiddleware,
		RequestMiddleware:   executor.requestMiddleware,
		Recover:             executor.recover,
		ErrorPresenter:      executor.handlerError,
		Tracer:              tracer,
	}
	defer func() {
		tracer.EndOperationExecution(ctx)
		err := reqCtx.RegisterExtension(tracingExt, tracer.GetTracing())
		if err != nil {
			executor.log.Error("failed register tracing extension", zap.Error(err))
		}

		resp.Extensions = reqCtx.Extensions
	}()

	ctx, doc, gqlErr := executor.parseQuery(ctx, reqCtx.Tracer, req.Query)
	if gqlErr != nil {
		resp = &graphql.Response{
			Errors: gqlerror.List{gqlErr},
		}
		return
	}
	reqCtx.Doc = doc

	ctx, op, vars, listErr := executor.validate(ctx, reqCtx.Tracer, reqCtx.Doc, req.OperationName, req.Variables)
	if len(listErr) != 0 {
		resp = &graphql.Response{
			Errors: listErr,
		}
		return
	}
	reqCtx.Variables = vars

	ctx = graphql.WithRequestContext(ctx, reqCtx)

	switch op.Operation {
	case ast.Query:
		resp = executor.Query(ctx, op)
	case ast.Mutation:
		resp = executor.Mutation(ctx, op)
	}

	return
}

func (executor *gqlGenExecutor) parseQuery(ctx context.Context, tracer graphql.Tracer, query string) (context.Context,
	*ast.QueryDocument, *gqlerror.Error) {
	ctx = tracer.StartOperationParsing(ctx)
	defer tracer.EndOperationParsing(ctx)

	doc, err := parser.ParseQuery(&ast.Source{Input: query})
	if err != nil {
		return ctx, nil, err
	}

	return ctx, doc, nil
}

func (executor *gqlGenExecutor) validate(ctx context.Context, tracer graphql.Tracer, doc *ast.QueryDocument, operationName string,
	variables map[string]interface{},
) (context.Context, *ast.OperationDefinition, map[string]interface{}, gqlerror.List) {
	ctx = tracer.StartOperationValidation(ctx)
	defer tracer.EndOperationValidation(ctx)

	schema := executor.Schema()

	listErr := validator.Validate(schema, doc)
	if len(listErr) != 0 {
		return ctx, nil, nil, listErr
	}

	op := doc.Operations.ForName(operationName)
	if op == nil {
		return ctx, nil, nil, gqlerror.List{gqlerror.Errorf("operation %s not found", operationName)}
	}

	vars, err := validator.VariableValues(schema, op, variables)
	if err != nil {
		return ctx, nil, nil, gqlerror.List{err}
	}

	return ctx, op, vars, nil
}

func (executor *gqlGenExecutor) requestMiddleware(ctx context.Context, next func(ctx context.Context) []byte) []byte {
	return next(ctx)
}

func (executor *gqlGenExecutor) handlerError(ctx context.Context, err error) (out *gqlerror.Error) {
	if err == nil {
		return
	}

	realErr := errors.Cause(err)

	if realErr == response.NotFoundError {
		return nil
	}

	if gqlErr, ok := realErr.(*gqlerror.Error); ok {
		return gqlErr
	}

	grpcErr, ok := status.FromError(realErr)
	if ok {
		grpcLog := executor.log.With(
			zap.Error(err),
			zap.Uint32("code", uint32(grpcErr.Code())),
			zap.Stringer("grpc_error", grpcErr.Proto()),
		)
		switch grpcErr.Code() {
		case codes.Canceled, codes.Unknown, codes.AlreadyExists, codes.PermissionDenied, codes.ResourceExhausted,
			codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Unimplemented, codes.Internal,
			codes.Unavailable:
			grpcLog.Error("grpc: internal error")
			return response.InternalError
		case codes.NotFound:
			grpcLog.Info("grpc: not found")
			return response.NotFoundError
		case codes.DeadlineExceeded:
			grpcLog.Info("grpc: timeout")
			return response.RequestTimeoutError
		case codes.InvalidArgument:
			grpcLog.Warn("grpc: bad arguments")
			return response.BadRequestError
		default:
			grpcCodeReplacing, ok := grpcCodeReplacing[grpcErr.Code()]
			if !ok {
				grpcLog.Warn("grpc: unknown code")
				return response.InternalError
			}
			return &gqlerror.Error{
				Message: grpcErr.Message(),
				Extensions: map[string]interface{}{
					response.CodeExt: grpcCodeReplacing,
				},
			}
		}
	}

	errReplacing, ok := errorsReplacing[err]
	if ok {
		return errReplacing
	}

	errContent, ok := contentErrorReplacing[err.Error()]
	if ok {
		return errContent
	}

	executor.log.Warn("unprocessed error", zap.Error(err))

	return response.InternalError
}

func (executor *gqlGenExecutor) recover(ctx context.Context, err interface{}) error {
	executor.log.Error("catch panic", zap.Reflect("error", err))
	return response.InternalError
}
