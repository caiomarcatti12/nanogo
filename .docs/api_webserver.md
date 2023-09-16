### WebServer

A biblioteca oferece uma maneira simplificada e estruturada para configurar e iniciar um servidor web, definir rotas e controladores para manipular requisições HTTP.

#### Configuração

Para configurar o servidor Web, você deve definir as seguintes variáveis de ambiente:

```sh
SERVER_PORT=8080
SERVER_MAX_UPLOAD_SIZE=5MB
```

- `SERVER_PORT`: A porta que será disponibilizada a aplicação.
- `SERVER_MAX_UPLOAD_SIZE`: Define o tamanho máximo permitido para upload por requisição. O valor padrão é 5MB.


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

### Manipulação de Uploads de Arquivos

Para manipular uploads de arquivos, você utilizará a estrutura `webserver.FileUpload`. Essa estrutura contém todas as informações necessárias sobre o arquivo que está sendo carregado. Abaixo, detalhamos cada campo da estrutura:

```go
type FileUpload struct {
	Filename string
	Size     int64
	Content  []byte
}
```

#### Campos

- **Filename**: Uma string que armazena o nome do arquivo enviado no upload.
- **Size**: Um int64 que representa o tamanho do arquivo em bytes.
- **Content**: Um slice de bytes (`[]byte`) que contém o conteúdo do arquivo enviado.

#### Uso no Controlador

Para utilizar o objeto `webserver.FileUpload` no seu controlador, você precisa incluí-lo na estrutura que representa a entrada de dados do controlador. Abaixo, temos um exemplo de como você pode fazer isso:

```go
package dto

type ApplicationCreateDTO struct {
	// ... outros campos
	FileData webserver.FileUpload   `form:"meu_arquivo"`
}
```

#### Exemplo de Controlador com Upload de Arquivo

Aqui está um exemplo de um controlador que manipula o upload de um arquivo:

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
	
	// Acessando os dados do arquivo carregado
	file := ctx.Payload.FileData
	
	//... sua lógica de negocio.

	return &webserver.APIResponse{
		Data:       response,
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
```


### Exemplo de Controlador com  Download de um Arquivo

Para facilitar o download de um arquivo a partir de um endpoint, você pode retornar o conteúdo do arquivo diretamente na resposta da API, configurando apropriadamente os cabeçalhos da resposta HTTP para indicar que um arquivo deve ser baixado. Aqui está um exemplo de como você pode configurar um handler para enviar um arquivo para download:

```go
package controller

import (
	"net/http"
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
)

func DownloadHandler(ctx *webserver.HandlerContext[dto.ApplicationCreateDTO]) (interface{}, error) {

	//... sua lógica de negocio.
	fileContent := []byte("Conteúdo do seu arquivo aqui")
	
	return &webserver.APIResponse{
		Data:       fileContent,
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":        "application/octet-stream",
			"Content-Disposition": "attachment; filename=meuarquivo.txt",
		},
	}, nil
}
```

Neste exemplo, o `DownloadHandler` está configurado para retornar um conjunto de bytes que representam o conteúdo do arquivo. Os cabeçalhos da resposta HTTP são configurados para indicar o tipo de conteúdo como "application/octet-stream", e um cabeçalho "Content-Disposition" está sendo usado para especificar que a resposta deve ser tratada como um arquivo para download, dando-lhe um nome de arquivo específico ("meuarquivo.txt").
