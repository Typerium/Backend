package graphql

import (
	"context"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func newTracer() *tracer {
	return &tracer{
		timing: make(map[string]interface{}),
	}
}

type tracer struct {
	mu     sync.Mutex
	timing map[string]interface{}
}

const (
	tracingExt = "tracing"

	tracingOperationParsing    = "parsing"
	tracingOperationValidation = "validation"

	tracingStart    = "startTime"
	tracingEnd      = "endTime"
	tracingDuration = "duration"
)

func (t *tracer) startTimer(ctx context.Context, key string) context.Context {
	t.mu.Lock()
	defer t.mu.Unlock()

	startTime, ok := t.timing[tracingStart].(time.Time)
	if !ok {
		return ctx
	}

	startTimer := time.Now().UTC()

	t.timing[key] = &timing{
		StartOffset: startTimer.Sub(startTime),
	}

	return context.WithValue(ctx, key, startTimer)
}

func (t *tracer) stopTimer(ctx context.Context, key string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	keyTiming, ok := t.timing[key].(*timing)
	if !ok {
		return
	}

	startTimer, ok := ctx.Value(key).(time.Time)
	if !ok {
		return
	}

	keyTiming.Duration = time.Now().UTC().Sub(startTimer)

	t.timing[key] = keyTiming
}

func (t *tracer) StartOperationParsing(ctx context.Context) context.Context {
	return t.startTimer(ctx, tracingOperationParsing)
}

func (t *tracer) EndOperationParsing(ctx context.Context) {
	t.stopTimer(ctx, tracingOperationParsing)
}

func (t *tracer) StartOperationValidation(ctx context.Context) context.Context {
	return t.startTimer(ctx, tracingOperationValidation)
}

func (t *tracer) EndOperationValidation(ctx context.Context) {
	t.stopTimer(ctx, tracingOperationValidation)
}

func (t *tracer) StartOperationExecution(ctx context.Context) context.Context {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now().UTC()
	t.timing[tracingStart] = now

	return context.WithValue(ctx, tracingStart, now)

}

func (t *tracer) EndOperationExecution(ctx context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	startTime, ok := t.timing[tracingStart].(time.Time)
	if !ok {
		return
	}

	now := time.Now().UTC()

	t.timing[tracingEnd] = now
	t.timing[tracingDuration] = now.Sub(startTime)
}

func (t *tracer) StartFieldExecution(ctx context.Context, field graphql.CollectedField) context.Context {
	return ctx
}

func (t *tracer) StartFieldResolverExecution(ctx context.Context, rc *graphql.ResolverContext) context.Context {
	return ctx
}

func (t *tracer) StartFieldChildExecution(ctx context.Context) context.Context {
	return ctx
}

func (t *tracer) EndFieldExecution(ctx context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.timing["execution"] = map[string]interface{}{
		"resolvers": []*resolver{},
	}
}

type timing struct {
	StartOffset time.Duration `json:"startOffset"`
	Duration    time.Duration `json:"duration"`
}

type resolver struct {
	Path       []string `json:"path"`
	ParentType string   `json:"parentType"`
	FieldName  string   `json:"fieldName"`
	ReturnType string   `json:"returnType"`
	timing
}

func (t *tracer) GetTracing() interface{} {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.timing["version"] = 1

	return t.timing
}
