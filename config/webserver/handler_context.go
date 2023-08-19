package webserver

import "net/http"

type HandlerContext struct {
	Payload  interface{}
	Headers  http.Header
	Request  *http.Request
	Response http.ResponseWriter
}