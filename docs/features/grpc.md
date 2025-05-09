# gRPC Web Server

Este pacote oferece uma estrutura simples e modular para a criação e gerenciamento de servidores gRPC em aplicações Go com suporte integrado a Dependency Injection (DI) e logs automáticos. Permite o registro dinâmico e modular de serviços, facilitando a expansão, manutenção e testabilidade.

## Estrutura

O pacote `grpc_webserver` inclui:

- **IGrpcServer:** Interface que representa o servidor gRPC.
- **Server:** Implementação da interface `IGrpcServer`, responsável pela inicialização e gestão do servidor.
- **Factory:** Função que cria instâncias do servidor gRPC com DI e logger injetados automaticamente.

## Funcionalidades

- Inicialização simplificada do servidor gRPC usando variáveis de ambiente.
- Registro modular de serviços através de handlers que utilizam automaticamente o sistema de DI.
- Logs informativos detalhados sobre o ciclo de vida do servidor e serviços.

## Uso Básico

### 1. Criando um Serviço

Cada serviço deve implementar uma interface compatível com a injeção de dependências:

```go
package example

import (
	"context"
	"google.golang.org/grpc"
	"your_package/pb"
)

type ExampleService struct {
	pb.UnimplementedExampleServer
}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (s *ExampleService) Register(server *grpc.Server) {
	pb.RegisterExampleServer(server, s)
}

// Handler exemplo do gRPC
func (s *ExampleService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Olá, " + req.Name}, nil
}
```

### 2. Inicializando o Servidor usando Factory

```go
package main

import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/env"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"your_package/grpc_webserver"
	"your_package/example"
)

func main() {
	container := di.GetInstance()
	environment := env.Factory(nil)
	grpcServer := grpc_webserver.Factory(container, environment)

	// Adicione seu serviço usando a estrutura GRPCHandler
	grpcServer.Add(grpc_webserver.GRPCHandler{
		IHandler:    example.NewExampleService,
		ServiceFunc: "Register",
	})

	// Inicie o servidor
	if err := grpcServer.Start(); err != nil {
		logger, _ := container.GetByFactory(log.Factory)
		logger.(log.ILog).Fatalf("Erro ao iniciar servidor gRPC: %v", err)
	}
}
```

### 3. Consumindo o Serviço

Exemplo simples de consumo do serviço usando um cliente gRPC:

```go
package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	"your_package/pb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewExampleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "João"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.Message) // Output: Greeting: Olá, João
}
```

## Variáveis de Ambiente

| Variável   | Descrição              | Default    |
|------------|------------------------|------------|
| GRPC_HOST  | Endereço do servidor   | `0.0.0.0`  |
| GRPC_PORT  | Porta do servidor gRPC | `50051`    |

## Métodos Principais

- `Factory(container di.IContainer, env env.IEnv) IGrpcServer`: Cria e configura uma nova instância do servidor.
- `Add(handler GRPCHandler)`: Registra um serviço no servidor.
- `Start() error`: Inicia o servidor no endereço configurado pelas variáveis de ambiente.

## Testes Automatizados

Execute testes automatizados com:

```shell
go test ./...
```

## Boas Práticas

- Sempre utilize o Factory para criar servidores gRPC.
- Configure os serviços utilizando DI para facilitar testes e manutenção.
- Utilize variáveis de ambiente claramente definidas para flexibilidade e escalabilidade.