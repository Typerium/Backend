package attributes

import (
	"context"
	"strings"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"typerium/internal/pkg/logging"
	"typerium/internal/pkg/web"
)

var (
	userAttrLog = logging.New("context_user_attributes")
)

func NewSchemaTokenAttribute() AttributeString {
	return &schemaTokenAttr{web.BearerSchema}
}

type schemaTokenAttr struct {
	schema string
}

func (attr *schemaTokenAttr) Get(ctx context.Context) (out string, ok bool) {
	out, ok = ctx.Value(tokenKey).(string)
	return
}

func (attr *schemaTokenAttr) Set(ctx context.Context, in string) context.Context {
	token := in[strings.Index(in, attr.schema)+len(attr.schema):]
	token = strings.TrimSpace(token)
	return context.WithValue(ctx, tokenKey, token)
}

func NewUserIDAttribute() AttributeUUID {
	return &userAttr{}
}

type userAttr struct{}

func (attr *userAttr) Get(ctx context.Context) (uuid.UUID, bool) {
	userIDCtx := ctx.Value(userIDKey)
	switch result := userIDCtx.(type) {
	case uuid.UUID:
		return result, true
	case string:
		out, err := uuid.FromString(result)
		if err != nil {
			userAttrLog.Warn("can't convert string to uuid", zap.Error(err))
			return uuid.Nil, false
		}
		return out, true
	}

	return uuid.Nil, false
}

func (attr *userAttr) Set(ctx context.Context, in uuid.UUID) context.Context {
	if uuid.Equal(in, uuid.Nil) {
		return ctx
	}
	return context.WithValue(ctx, userIDKey, in)
}
