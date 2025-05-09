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
package grcp_webserver

import (
	"net"

	"google.golang.org/grpc"
)

// GRPCService define a interface que todos os serviços gRPC devem implementar.
type GRPCService interface {
	// Register registra o serviço no servidor gRPC fornecido.
	Register(server *grpc.Server)
}

// Server encapsula a inicialização e o gerenciamento de um servidor gRPC,
// permitindo a adição modular de serviços que implementam a interface GRPCService.
type Server struct {
	grpc     *grpc.Server  // Instância do servidor gRPC.
	services []GRPCService // Serviços a serem registrados no servidor.
}

// New cria e retorna uma nova instância do Server.
// Aceita opções variáveis que serão passadas diretamente ao servidor gRPC interno.
func New(opts ...grpc.ServerOption) *Server {
	return &Server{
		grpc: grpc.NewServer(opts...),
	}
}

// Add adiciona um serviço ao servidor.
// O serviço deve implementar a interface GRPCService.
func (s *Server) Add(service GRPCService) {
	s.services = append(s.services, service)
}

// Start inicia o servidor gRPC escutando no endereço fornecido.
// Este método realiza as seguintes etapas:
// - Cria um listener TCP no endereço especificado.
// - Registra todos os serviços adicionados previamente ao servidor gRPC.
// - Inicia o servidor gRPC e permanece bloqueado enquanto o servidor estiver em execução.
// Retorna um erro caso não seja possível escutar no endereço fornecido ou o servidor não possa iniciar.
func (s *Server) Start(address string) error {
	// Criando um listener TCP para o endereço fornecido
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Registrando todos os serviços adicionados ao servidor
	for _, svc := range s.services {
		svc.Register(s.grpc)
	}

	// Iniciando o servidor gRPC
	return s.grpc.Serve(lis)
}
