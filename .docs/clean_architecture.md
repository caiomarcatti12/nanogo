# Framework nanogo e Clean Architecture

Neste segmento da arquitetura, focamos na interação com sistemas e ferramentas externas, como bancos de dados, sistemas de arquivos, entre outros. Esta camada contém os detalhes de como as operações externas são realizadas.

O framework **nanogo** destaca-se exatamente por facilitar a implementação desta camada, oferecendo soluções eficientes e descomplicadas para integrações com diferentes sistemas externos, garantindo assim, uma base sólida para o desenvolvimento da sua aplicação.


# Clean Architecture

A **Arquitetura Clean** ou "Clean Architecture" foi proposta por Robert C. Martin (também conhecido como Uncle Bob), e tem como objetivo principal a separação de preocupações, tornando os sistemas mais organizados e facilitando a manutenção e a expansão do código.

Aqui, apresentamos um guia sobre como podemos implementar e beneficiar-se da Arquitetura Clean em projetos Go.

## Benefícios

- **Manutenção facilitada**: Devido à separação clara e definida de cada componente, a manutenção torna-se mais simples e rápida.
- **Testabilidade**: Facilita a escrita e manutenção de testes unitários.
- **Independência de Frameworks**: O código do núcleo da aplicação fica independente de qualquer framework externo.
- **Flexibilidade**: Permite mudanças significativas na aplicação com menor esforço.
- **Desacoplamento**: Facilita a substituição de componentes sem afetar outras partes do sistema.

## Estrutura Principal

A Arquitetura Clean pode ser dividida em camadas, cada uma com uma responsabilidade específica:

1. **Entities**: Contém a lógica central do negócio e as regras que não devem ser afetadas por mudanças externas.
2. **Use Cases**: Define as operações disponíveis no sistema do ponto de vista do usuário final.
3. **Controllers/Adapters**: Atua como um conversor, onde os dados de entrada e saída são convertidos para um formato que pode ser entendido pela próxima camada.
4. **External Interfaces (Framework & Drivers)**: Contém detalhes sobre os frameworks e ferramentas utilizados no sistema (como bancos de dados, servidores web, etc).

### Entities

Nesta camada, definimos as estruturas e métodos que representam e operam sobre os objetos de negócio centrais na nossa aplicação. Por exemplo, uma entidade "Usuario" em uma aplicação poderia ser representada como:

```go
package entity

type Usuario struct {
	ID    string
	Nome  string
	Email string
}

func (u *Usuario) ValidarEmail() error {
	// Lógica de validação de e-mail
	return nil
}
```

### Use Cases

Aqui definimos as ações que podem ser executadas em nossas entidades. Estas ações representam a lógica de negócios da aplicação.

```go
package usecase

import (
	"meuprojeto/entity"
)

func CriarUsuario(usuario entity.Usuario) error {
	// Lógica para criar um novo usuário
	return nil
}
```

### Controllers/Adapters

Esta camada é responsável por receber as requisições externas, passando-as para a camada de Use Cases após adequada conversão/validação. Veja um exemplo de um controlador HTTP que utiliza o pacote `net/http` para definir rotas:

```go
package controller

import (
	"meuprojeto/usecase"
	"net/http"
)

func CriarUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	// Lógica para lidar com a requisição HTTP e chamar o use case correspondente
}
```

### External Interfaces (Framework & Drivers)

Aqui, temos o código que irá interagir com serviços externos, como bancos de dados, sistemas de arquivos, entre outros. Por exemplo, uma implementação de um repositório que utiliza um banco de dados para persistir dados de usuários:

```go
package repository

import "meuprojeto/entity"

type UsuarioRepository struct {
	//...
}

func (repo *UsuarioRepository) Criar(usuario entity.Usuario) error {
	// Lógica para criar um usuário no banco de dados
	return nil
}
```

