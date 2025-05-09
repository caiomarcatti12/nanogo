package grpc_webserver

type GRPCHandler struct {
	IHandler    interface{}
	ServiceFunc string
}
