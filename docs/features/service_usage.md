# Estrutura e Uso de Services no Nanogo

O padrão Service é fundamental para organizar a lógica de negócio no Nanogo, promovendo separação de responsabilidades, testabilidade e integração fácil com o sistema de Dependency Injection (DI).

## O que é um Service?
Um Service encapsula regras de negócio e operações de domínio, sendo consumido por controllers, handlers ou outros serviços. Ele não deve conter lógica de infraestrutura (ex: acesso direto a banco, HTTP, etc), mas pode depender de repositórios ou gateways.

## Estrutura Básica de um Service

O ideal é sempre definir uma interface para o seu service, facilitando a inversão de dependência, testes e mocks. O construtor (New...) deve retornar a interface, não a struct concreta.

```go
package example

type IExampleService interface {
    DoSomething(param string) (string, error)
}

type ExampleService struct {
    repo example.IExampleRepository // dependência injetada
}

func NewExampleService(repo example.IExampleRepository) IExampleService {
    return &ExampleService{repo: repo}
}

func (s *ExampleService) DoSomething(param string) (string, error) {
    // lógica de negócio
    result, err := s.repo.FindByParam(param)
    if err != nil {
        return "", err
    }
    return "Resultado: " + result, nil
}
```

## Registrando um Service no DI
O registro é feito no bootstrap do projeto:

```go
import "github.com/caiomarcatti12/nanogo/pkg/di"

di.RegisterFactory(NewExampleService)
```

## Consumindo um Service
Controllers, handlers ou outros serviços podem obter instâncias via DI:

```go
import "github.com/caiomarcatti12/nanogo/pkg/di"

service, _ := di.Get[IExampleService]()
```

## Boas Práticas
- Sempre injete dependências via construtor.
- Sempre exponha uma interface para o service.
- Não acople lógica de infraestrutura diretamente no service.
- Escreva interfaces para facilitar testes e mocks.
- Mantenha métodos pequenos e focados em uma responsabilidade.

---

> Consulte também: [Arquitetura Limpa](./clean_architecture.md) e [Repository Pattern](./repository_architecture.md)
