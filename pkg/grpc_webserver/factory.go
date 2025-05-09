package grpc_webserver

import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"google.golang.org/grpc"
)

// IGrpcServer é a interface que representa nosso encapsulamento do servidor gRPC.
type IGrpcServer interface {
	Add(handler GRPCHandler)
	Start() error
}

// Factory cria uma instância do servidor gRPC com DI e logger injetados automaticamente.
func Factory(logger log.ILog, env env.IEnv) IGrpcServer {
	host := env.GetEnv("GRPC_HOST", "0.0.0.0")
	port := env.GetEnv("GRPC_PORT", "50051")

	interceptors := grpc.ChainUnaryInterceptor(
		correlationIdInterceptor(),
	)

	return &Server{
		grpc: grpc.NewServer(interceptors),

		handlers: []GRPCHandler{},
		di:       di.GetInstance(),
		logger:   logger.(log.ILog),
		host:     host,
		port:     port,
	}
}
