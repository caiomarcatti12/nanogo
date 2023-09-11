### WebServer

A biblioteca oferece uma maneira simplificada e estruturada para configurar e iniciar um servidor web, definir rotas e controladores para manipular requisições HTTP.

#### Configuração

Para configurar o servidor Web, você deve definir as seguintes variáveis de ambiente:

```sh
SERVER_PORT=8080
```

- `SERVER_PORT`: A porta que será disponibilizada a aplicação.


#### Inicialização do Servidor Web

Aqui está um exemplo de como inicializar o servidor web:

```go
package main

import (
	"github.com/caiomarcatti12/minha-aplicacao/controller"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
)

func main() {
	env.LoadEnv()

	server := webserver.NewWebServer()

	server.Start()
}
```
#### Inicialização do Servidor Web HTTPS (TLS)

Para configurar o HTTPS no servidor web HTTPS, siga a documentação - **[Configuração HTTPS (TLS)](./api_webserver_tls.md)**

#### Definindo Rotas

As rotas são definidas com uso do pacote `webserver`, onde você pode especificar o método HTTP, o caminho da URL, o manipulador e a estrutura DTO para a validação da carga útil, conforme mostrado abaixo:

```go
package controller

import (
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
)

func ApplicationRouter() {
	webserver.AddRouter("DELETE", "/{id}", DeleteHandler, dto.IDStruct{})
	webserver.AddRouter("PUT", "/{id}", UpdateHandler, dto.ApplicationUpdateDTO{})
	webserver.AddRouter("GET", "/{id}", FindByIdHandler, dto.IDStruct{})
	webserver.AddRouter("GET", "/", FindAllHandler)
	webserver.AddRouter("POST", "/", CreateHandler, dto.ApplicationCreateDTO{})
}
```

#### Criando um Controlador

Aqui está um exemplo de como criar um controlador para manipular requisições HTTP POST para criar uma nova aplicação:

```go
package controller

import (
	"github.com/caiomarcatti12/minha-aplicacao/dto"
	"net/http"
	"github.com/caiomarcatti12/minha-aplicacao/repository"
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
	"github.com/caiomarcatti12/minha-aplicacao/service"
)

func CreateHandler(ctx *webserver.HandlerContext[dto.ApplicationCreateDTO]) (interface{}, error) {
	
	//... sua lógica de negocio.

	return &webserver.APIResponse{
		Data:       response,
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
```

### Explicando o Método CreateHandler

O método `CreateHandler` é responsável por lidar com as requisições HTTP POST que visam criar uma nova aplicação. Abaixo, detalhamos os parâmetros de entrada e saída.

#### Parâmetros de Entrada

- **ctx *webserver.HandlerContext[dto.ApplicationCreateDTO]**: Este é o contexto do manipulador que contém todas as informações sobre a requisição HTTP atual, incluindo o payload. O contexto é tipado com a estrutura `dto.ApplicationCreateDTO` para garantir que o payload atenda ao formato esperado.

- **ctx.Payload**: O payload contém os dados enviados na requisição HTTP. Estes dados são deserializados na estrutura `dto.ApplicationCreateDTO`.

#### Saída

- **interface{}**: Este é o tipo de retorno que pode conter diferentes tipos de estruturas de resposta. No caso, está retornando uma estrutura `webserver.APIResponse`.

- **error**: Se ocorrer um erro durante o processamento, ele será retornado aqui.

A estrutura `webserver.APIResponse` contém os seguintes campos:

- **Data**: Contém os dados da resposta, que neste caso, são os detalhes da aplicação criada.

- **StatusCode**: O código de status HTTP para a resposta. Aqui é utilizado o código 201, que indica que um novo recurso foi criado com sucesso.

- **Headers**: Um mapa contendo os cabeçalhos HTTP da resposta. Neste exemplo, está sendo definido o cabeçalho "Content-Type" como "application/json".