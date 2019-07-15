package response

const (
	// CodeExt key for extension of error
	CodeExt = "code"
)

// values of CodeExt key
const (
	InternalErrorCode  = "500"
	NotFoundErrorCode  = "404"
	BadRequestCode     = "400"
	RequestTimeoutCode = "408"
	UnauthorizedCode   = "401"
	ForbiddenCode      = "403"
)
