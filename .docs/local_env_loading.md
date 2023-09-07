# **Carregamento Local de Variáveis de Ambiente**

O Nanogo facilita o gerenciamento de suas variáveis de ambiente através do carregamento local de um arquivo `.env` situado na raiz do seu projeto. Este recurso é essencial para manter uma configuração segura e organizada.

## **LoadEnv**

A função `LoadEnv` é responsável por carregar as variáveis de ambiente a partir deste arquivo `.env`. O método lê e define as variáveis de ambiente conforme especificado no arquivo. Se houver algum problema durante o carregamento, uma mensagem de erro será registrada e o processo será encerrado.

**Exemplo de uso:**

```go
func main() {
	env.LoadEnv()
}
```

## **Exemplo de arquivo `.env`**

Aqui está um exemplo de um arquivo `.env` para iniciar:

```env
APP_NAME=MeuAppIncrivel
ENV=producao
VAR_1=VALUE_1
```

Utilizando este arquivo, você pode manter as configurações essenciais do seu aplicativo de forma estruturada e segura.