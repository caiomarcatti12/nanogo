package webserver

import "net/http"

type HandlerContext[T any] struct {
	Payload  T
	RawQuery string
	Headers  http.Header
	Request  *http.Request
	Response http.ResponseWriter
}
