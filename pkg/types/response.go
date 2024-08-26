package types

type Response struct {
	Data       interface{}
	StatusCode int
	Headers    map[string]string
}
