### Configuração HTTPS (TLS)

Para configurar o HTTPS no servidor web, é necessário especificar os caminhos para o certificado e a chave privada através das variáveis de ambiente `SERVER_CERTIFICATE` e `SERVER_KEY`. Essa configuração habilita o TLS (Transport Layer Security), garantindo que as comunicações entre o servidor e os clientes sejam criptografadas e seguras.

#### Variáveis de Ambiente

```
SERVER_CERTIFICATE=./server/server.crt
SERVER_KEY=./server/server.key
```

- `SERVER_CERTIFICATE`: O caminho para o arquivo do certificado do servidor (geralmente com a extensão `.crt` ou `.pem`).
- `SERVER_KEY`: O caminho para o arquivo da chave privada do servidor (geralmente com a extensão `.key`).

Certifique-se de que os caminhos especificados estão corretos e que os arquivos estão acessíveis pelo servidor para evitar problemas de inicialização.

#### Geração de Certificado e Chave Autoassinados

Caso não possua um certificado emitido por uma Autoridade Certificadora (CA), você pode criar um certificado e uma chave autoassinados seguindo os passos abaixo. Isso pode ser útil para ambientes de teste e desenvolvimento:

1. **Gere uma chave privada RSA**:
   ```sh
   openssl genpkey -algorithm RSA -out server.key
   ```
   Este comando cria uma chave privada RSA, que será usada para criar o certificado.

2. **Crie uma solicitação de assinatura de certificado (CSR)**:
   ```sh
   openssl req -new -key server.key -out server.csr -subj "/C=US/ST=California/L=San Francisco/O=My Organization/OU=My Organizational Unit/CN=host.docker.internal/emailAddress=myemail@mydomain.com"
   ```
   Este comando cria uma CSR contendo os detalhes da sua organização. A CSR é usada para criar o certificado autoassinado.

3. **Crie um certificado autoassinado**:
   ```sh
   openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
   ```
   Este comando cria um certificado autoassinado válido por 365 dias usando a CSR e a chave privada geradas nos passos anteriores.

#### Integrando TLS no Servidor

Após configurar as variáveis de ambiente com os caminhos corretos para o certificado e a chave, você precisará ajustar a inicialização do servidor web para usar HTTPS. Aqui está um exemplo genérico:

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

Neste exemplo, `Start` é chamado e internamente ele inicia o servidor HTTPS da mesma maneira que o servidor HTTP

#### Observações

- **Segurança**: Certificados autoassinados podem gerar avisos de segurança nos navegadores, pois eles não foram emitidos por uma CA confiável. Para ambientes de produção, é altamente recomendado obter um certificado de uma CA confiável.
- **Renovação**: Lembre-se de renovar o certificado antes que expire, para evitar interrupções no serviço.
- **Informações de Certificado**: O comando `openssl req` inclui uma opção `-subj` para especificar as informações de certificado. Adapte esses valores conforme necessário para corresponder às informações da sua organização.

