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
func Factory(container di.IContainer, env env.IEnv) IGrpcServer {
	logger, err := container.GetByFactory(log.Factory)
	if err != nil {
		panic("Não foi possível obter o logger: " + err.Error())
	}

	host := env.GetEnv("GRPC_HOST", "0.0.0.0")
	port := env.GetEnv("GRPC_PORT", "50051")

	logger.(log.ILog).Infof("Iniciando servidor gRPC em %s:%s", host, port)

	return &Server{
		grpc:     grpc.NewServer(),
		handlers: []GRPCHandler{},
		di:       container,
		logger:   logger.(log.ILog),
		host:     host,
		port:     port,
	}
}
