package webserver

type APIResponse struct {
	Data       interface{}
	StatusCode int
	Headers    map[string]string
}