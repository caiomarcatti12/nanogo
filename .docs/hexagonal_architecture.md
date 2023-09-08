# Framework nanogo e Arquitetura Hexagonal

A arquitetura hexagonal, também conhecida como "Ports and Adapters", procura manter o núcleo da aplicação (a lógica de negócios) isolada de qualquer influência externa, como bancos de dados, sistemas de arquivos e interfaces de usuário. No âmbito deste sistema arquitetônico, o **nanogo** surge como uma ferramenta vital, facilitando a implementação ao encapsular bibliotecas externas, promovendo uma série de benefícios significativos que detalhamos abaixo.

## Entendendo a Arquitetura Hexagonal

A estrutura da arquitetura hexagonal pode ser visualizada como uma aplicação hexagonal, onde cada lado representa um port, que podem ser divididos em "Ports Primários" e "Ports Secundários", intermediados por "Adapters".

### Encapsulamento das Bibliotecas Externas com nanogo

O **nanogo** proporciona uma maneira eficiente e organizada de encapsular bibliotecas externas, criando uma abstração sólida e facilitando a integração destas bibliotecas no seu projeto. Este encapsulamento tem vários benefícios, incluindo:

- **Simplicidade**: Ao abstrair os detalhes complexos das bibliotecas externas, o nanogo permite um desenvolvimento mais rápido e direto, facilitando a utilização destas bibliotecas.
- **Segurança**: O encapsulamento ajuda a criar uma barreira de segurança, onde os detalhes internos das bibliotecas estão seguros e protegidos de acessos indevidos.
- **Manutenção**: Alterações ou atualizações nas bibliotecas externas podem ser feitas de forma centralizada, sem afetar diretamente os consumidores destas bibliotecas.
- **Testabilidade**: Facilita a criação de mocks para testes unitários, garantindo que você pode testar sua lógica de negócios de forma isolada, sem a necessidade de integrar com serviços externos.

#### Exemplo de Encapsulamento

No contexto do **nanogo**, o encapsulamento pode ser observado na forma como gerenciamos e configuramos serviços externos, como na inicialização do servidor web:

```go
package main

import (
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
)

func main() {
	env.LoadEnv()
	server := webserver.NewWebServer()
	server.Start()
}
```

Aqui, a complexidade da configuração e inicialização do servidor web é encapsulada, proporcionando uma inicialização simplificada e focada.

### Implementação dos Ports e Adapters com nanogo

Com o auxílio do **nanogo**, a implementação dos ports e adapters torna-se uma tarefa menos árdua, permitindo uma clara delimitação entre a lógica de negócios e os agentes externos, conforme detalhado nas seções de "Ports Primários e Secundários" e "Adapters".

### Conclusão

A arquitetura hexagonal, auxiliada pelo framework **nanogo**, promove uma estruturação de projeto que beneficia a manutenção, segurança, e desenvolvimento. A prática de encapsular bibliotecas externas, inerente ao **nanogo**, não apenas protege o núcleo da sua aplicação, mas também proporciona uma experiência de desenvolvimento mais eficiente e organizada.
