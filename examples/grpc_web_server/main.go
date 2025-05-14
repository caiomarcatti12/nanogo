package main

import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/grpc_webserver"
	"github.com/caiomarcatti12/nanogo/pkg/nanogo"
)

func main() {
	nanogo.Bootstrap()

	grpcServer, err := di.Get[grpc_webserver.IGrpcServer]()

	if err != nil {
		panic("não foi possível obter o logger: " + err.Error())
	}

	// Adicione seu handler usando DI
	grpcServer.Add(grpc_webserver.GRPCHandler{
		IHandler:    NewEventConsumerHandler,
		ServiceFunc: "Register",
	})

	if err := grpcServer.Start(); err != nil {
	}
}
