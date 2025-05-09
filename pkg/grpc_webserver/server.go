/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package grpc_webserver

import (
	"fmt"
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"net"
	"reflect"

	"google.golang.org/grpc"
)

// Server implementa IGrpcServer.
type Server struct {
	grpc     *grpc.Server
	handlers []GRPCHandler
	di       di.IContainer
	logger   log.ILog
	host     string
	port     string
}

func (s *Server) Add(handler GRPCHandler) {
	err := s.di.Register(handler.IHandler)
	if err != nil {
		s.logger.Errorf("Falha ao registrar o handler: %v", err)
	}
	s.handlers = append(s.handlers, handler)
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.host, s.port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Errorf("Falha ao escutar em %s: %v", address, err)
		return err
	}

	for _, handler := range s.handlers {
		instance, err := s.di.GetByFactory(handler.IHandler)
		if err != nil {
			s.logger.Errorf("Falha ao obter instância do handler: %v", err)
			return err
		}

		method := reflect.ValueOf(instance).MethodByName(handler.ServiceFunc)
		if !method.IsValid() {
			return fmt.Errorf("método de registro de serviço não encontrado: %s", handler.ServiceFunc)
		}

		method.Call([]reflect.Value{reflect.ValueOf(s.grpc)})
	}

	s.logger.Infof("Servidor gRPC iniciado com sucesso em %s", address)
	return s.grpc.Serve(lis)
}
