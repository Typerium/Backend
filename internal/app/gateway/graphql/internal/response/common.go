package response

import (
	"github.com/vektah/gqlparser/gqlerror"
)

// graphql errors
var (
	InternalError = &gqlerror.Error{
		Message: "internal error",
		Extensions: map[string]interface{}{
			CodeExt: InternalErrorCode,
		},
	}
	NotFoundError = &gqlerror.Error{
		Message: "not found",
		Extensions: map[string]interface{}{
			CodeExt: NotFoundErrorCode,
		},
	}
	BadRequestError = &gqlerror.Error{
		Message: "bad request",
		Extensions: map[string]interface{}{
			CodeExt: BadRequestCode,
		},
	}
	RequestTimeoutError = &gqlerror.Error{
		Message: "request timeout",
		Extensions: map[string]interface{}{
			CodeExt: RequestTimeoutCode,
		},
	}
	Unauthorized = &gqlerror.Error{
		Message: "unauthorized",
		Extensions: map[string]interface{}{
			CodeExt: UnauthorizedCode,
		},
	}
	Forbidden = &gqlerror.Error{
		Message: "forbidden",
		Extensions: map[string]interface{}{
			CodeExt: ForbiddenCode,
		},
	}
)
