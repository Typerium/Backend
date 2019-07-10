package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

//go:generate go run github.com/99designs/gqlgen --config gqlgen.yml

type Request struct {
	Query         string
	OperationName string
	Variables     map[string]interface{}
}

type Executor interface {
	Exec(ctx context.Context, request *Request) []byte
}

func New() Executor {
	return nil
}

type gqlGenExecutor struct {
	graphql.ExecutableSchema
}

func (executor *gqlGenExecutor) Exec(ctx context.Context, request *Request) []byte {
	panic("implement me")
}
