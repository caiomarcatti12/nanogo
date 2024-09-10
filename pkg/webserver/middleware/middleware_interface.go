package webserver_middleware

import "net/http"

type IMiddleware interface {
	GetName() string
	Process(w http.ResponseWriter, r *http.Request, next http.Handler)
}
