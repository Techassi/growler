package response

import (
	"net/http"
)

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}
