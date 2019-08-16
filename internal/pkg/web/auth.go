package web

import (
	"encoding/base64"
	"fmt"
)

var (
	basicAuthPrefix = []byte("Basic ")
)

func BasicAuth(username, password string) []byte {
	src := []byte(fmt.Sprintf("%s:%s", username, password))
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))

	base64.StdEncoding.Encode(buf, src)

	return append(basicAuthPrefix, buf ...)
}

const (
	BearerSchema = "Bearer"
)