package web

import (
	"fmt"
)

type httpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *httpError) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"code": %d, "msg": "%s"}`, e.Code, e.Msg)), nil
}

// http errors
var (
	InternalServerErr = &httpError{
		Code: 500,
		Msg:  "internal server error",
	}
	NotFoundErr = &httpError{
		Code: 404,
		Msg:  "not found",
	}
	MethodNotAllowedErr = &httpError{
		Code: 405,
		Msg:  "method not allowed",
	}
)
