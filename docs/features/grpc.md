# gRPC Web Server

Este pacote oferece uma estrutura simples e modular para a criação e gerenciamento de servidores gRPC em aplicações Go. Permite o registro dinâmico e modular de serviços, facilitando a expansão e manutenção.

## Estrutura

O pacote `grpc-webserver` inclui:

* **Server:** Classe principal responsável pela inicialização e gestão do servidor gRPC.
* **GRPCService:** Interface que deve ser implementada por todos os serviços gRPC a serem registrados.

## Funcionalidades

* Inicialização simples do servidor gRPC.
* Registro modular de serviços através da interface `GRPCService`.
* Tratamento automático de conexões TCP.

## Uso Básico

### 1. Criando um Serviço

Cada serviço deve implementar a interface `GRPCService`:

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

func (s *ExampleService) Register(server *grpc.Server) {
	pb.RegisterExampleServer(server, s)
}

// Handler exemplo do gRPC
func (s *ExampleService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Olá, " + req.Name}, nil
}
```

### 2. Inicializando o Servidor

```go
package main

import (
	"log"
	"your_package/grpc_webserver"
	"your_package/example"
)

func main() {
	server := grpc_webserver.New()

	// Adicione seu serviço
	server.Add(&example.ExampleService{})

	// Inicie o servidor
	if err := server.Start(":50051"); err != nil {
		log.Fatalf("failed to start server: %v", err)
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

## Métodos Principais

* `New(opts ...grpc.ServerOption) *Server`: Cria e inicializa uma nova instância do servidor.
* `Add(service GRPCService)`: Adiciona um serviço ao servidor.
* `Start(address string) error`: Inicia o servidor na porta TCP especificada.

## Tratamento de Erros

O servidor retorna erros diretamente ao iniciar, permitindo que o chamador gerencie falhas na criação do listener ou ao iniciar o servidor gRPC.

## Testes Automatizados

Para garantir a estabilidade e robustez do pacote, testes são executados usando:

```shell
go test ./...
```

## Boas Práticas

* Sempre valide parâmetros antes de adicionar serviços ao servidor.
* Utilize portas claramente definidas e documentadas para evitar conflitos.
* Registre todos os serviços antes de iniciar o servidor para garantir que estejam disponíveis imediatamente após o início do servidor.
