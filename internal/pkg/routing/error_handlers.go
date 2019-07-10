package routing

import (
	"encoding/json"
	"fmt"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type httpError struct {
	Code    int
	Message string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("code %d; message: '%s'", e.Code, e.Message)
}

func (e *httpError) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{ \"code\": %d, \"message\": \"%s\" }", e.Code, e.Message)), nil
}

var (
	notFoundError = &httpError{
		Code:    fasthttp.StatusNotFound,
		Message: "not found",
	}
	internalError = &httpError{
		Code:    fasthttp.StatusInternalServerError,
		Message: "internal error",
	}
	methodNotAllowedError = &httpError{
		Code:    fasthttp.StatusMethodNotAllowed,
		Message: "method not allowed",
	}
)

func errorHandler(resp *fasthttp.Response, e *httpError) {
	resp.SetStatusCode(e.Code)
	data, err := json.Marshal(e)
	if err != nil {
		log.Error(errors.Wrap(err, "can't marshal json"))
		return
	}
	resp.SetBody(data)
}
