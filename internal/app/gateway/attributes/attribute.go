package attributes

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type AttributeString interface {
	Get(ctx context.Context) (out string, ok bool)
	Set(ctx context.Context, in string) context.Context
}

type AttributeUUID interface {
	Get(ctx context.Context) (out uuid.UUID, ok bool)
	Set(ctx context.Context, in uuid.UUID) context.Context
}
