package webserver

import "net/http"

type HandlerContext[T any] struct {
	Payload  T
	Headers  http.Header
	Request  *http.Request
	Response http.ResponseWriter
}
