# **Carregamento Remoto de Variáveis de Ambiente**

Além do carregamento local, o Nanogo oferece uma solução dinâmica e centralizada para gerenciar as configurações do seu aplicativo através do carregamento remoto de variáveis de ambiente.

## **Usando em Conjunto com Arquivo `.env` Local**

Ao trabalhar com o carregamento remoto de variáveis de ambiente, é recomendável definir as configurações básicas necessárias para conectar ao servidor remoto no seu arquivo `.env` local. Dessa forma, as informações críticas como o host do servidor, o token de autenticação, o nome do aplicativo e o ambiente estão prontamente disponíveis ao inicializar seu aplicativo.

Aqui está um exemplo de como as variáveis básicas podem ser definidas em seu arquivo `.env`:

```env
CLOUD_PROPERTIES_HOST=https://seu-servidor.com
CLOUD_PROPERTIES_TOKEN=seu_token_secreto
APP_NAME=MeuAppIncrivel
ENV=producao
ENV_REFRESH_TIME=5
```

Depois de definir essas variáveis básicas no arquivo `.env` local, você pode utilizá-las para configurar os parâmetros de `LoadRemoteEnvParams` e iniciar o carregamento remoto de variáveis de ambiente:

Os parâmetros para a função `LoadRemoteEnv` são opcionais. Caso você não forneça os valores durante a chamada, eles serão buscados a partir das respectivas variáveis de ambiente configuradas em seu arquivo `.env` ou no ambiente de execução atual.

**Exemplo de uso:**

```go
func main() {
	params := env.LoadRemoteEnvParams{
		Host:     env.GetEnv("CLOUD_PROPERTIES_HOST"),
		Token:    env.GetEnv("CLOUD_PROPERTIES_TOKEN"),
		AppName:  env.GetEnv("APP_NAME"),
		Env:      env.GetEnv("ENV"),
		Attempts: 1,
	}
	
	env.LoadRemoteEnv(params)
	
	// Os parâmetros para a função `LoadRemoteEnv` são opcionais.
	env.LoadRemoteEnv()
}
```

## **LoadRemoteEnv**

A função `LoadRemoteEnv` permite buscar variáveis de ambiente de um servidor remoto. Você deve fornecer uma estrutura `LoadRemoteEnvParams` com detalhes como o host, token de autenticação, nome do aplicativo, o ambiente atual e o número de tentativas para a busca das configurações.

A resposta do servidor deve estar em um formato JSON específico, contendo pares de chave-valor para cada variável de ambiente:

```json
{
    "CHAVE1": "VALOR1",
    "CHAVE2": "VALOR2"
    // ...
}
```

A função irá iterar sobre cada par, configurando as variáveis de ambiente em seu aplicativo. Existe também uma funcionalidade de auto-atualização que renova as variáveis de ambiente em intervalos determinados.

### Verifique o Carregamento Dinâmico de Variáveis de Ambiente

Para verificar manualmente o funcionamento do carregamento dinâmico de variáveis de ambiente, você pode utilizar o comando `curl` para simular a request que o `nanogo` faz ao iniciar. Este teste pode ser essencial durante a fase de desenvolvimento para garantir que o servidor está respondendo corretamente e fornecendo as configurações necessárias.

Antes de realizar o teste, assegure-se de que as variáveis `CLOUD_PROPERTIES_HOST`, `APP_NAME`, `ENV`, e `CLOUD_PROPERTIES_TOKEN` estão corretamente configuradas no seu ambiente ou substitua-as diretamente no comando com os valores apropriados. Aqui está um exemplo de como você pode fazer isso:

```sh
curl -X GET \
     -H "Accept: application/json" \
     -H "Authorization: Bearer ${CLOUD_PROPERTIES_TOKEN}" \
     "${CLOUD_PROPERTIES_HOST}/${APP_NAME}/${ENV}"
```

Este comando realiza uma solicitação GET ao servidor configurado, com os detalhes necessários para recuperar as variáveis de ambiente. Note que o token de autorização é passado no cabeçalho para autenticar a request.

**Nota:** Caso o `CLOUD_PROPERTIES_TOKEN` não seja necessário, você pode omitir a linha correspondente ao cabeçalho "Authorization".
