package routing

// value of Content-Type header
const (
	TextContentType = "text/plain; charset=utf-8"
	JSONContentType = "application/json"
)

const (
	xContentTypeOptionsNosniff = "nosniff"
)

const (
	acceptHeader        = "Accept"
	contentTypeHeader   = "Content-Type"
	authorizationHeader = "Authorization"

	originHeader                    = "Origin"
	accessControlAllowOriginHeader  = "Access-Control-Allow-Origin"
	accessControlAllowMethods       = "Access-Control-Allow-Methods"
	accessControlAllowHeadersHeader = "Access-Control-Allow-Headers"

	requestHeader             = "X-Request-ID"
	xContentTypeOptionsHeader = "X-Content-Type-Options"
)
