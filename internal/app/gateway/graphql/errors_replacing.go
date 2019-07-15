package graphql

import (
	"github.com/vektah/gqlparser/gqlerror"
	"google.golang.org/grpc/codes"

	"typerium/internal/app/gateway/graphql/internal/response"
	"typerium/internal/app/gateway/graphql/internal/types"
)

var (
	grpcCodeReplacing = map[codes.Code]string{
	}
)

var (
	errorsReplacing = map[error]*gqlerror.Error{
		types.ErrWrongType: {
			Message: types.ErrWrongType.Error(),
			Extensions: map[string]interface{}{
				response.CodeExt: response.BadRequestCode,
			},
		},
	}
)

var (
	contentErrorReplacing = map[string]*gqlerror.Error{
		"wrong type": {
			Message: "wrong type",
			Extensions: map[string]interface{}{
				response.CodeExt: response.BadRequestCode,
			},
		},
	}
)
