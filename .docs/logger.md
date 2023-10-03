### Logger

A classe `Logger` fornece uma maneira centralizada e estruturada de gerenciar mensagens de log em sua aplicação. Ela é integrada com configurações de ambiente e pode gerar logs em diferentes formatos com base no ambiente atual.

Claro, vou ajustar a seção sobre o Correlation ID para incluir essas informações:

#### Correlation ID

O logger usa um Correlation ID para rastrear e correlacionar logs originados a partir de uma única operação ou requisição. Este identificador ajuda a agrupar e identificar logs associados a uma ação específica, facilitando a análise e a depuração.

**Origem do Correlation ID**:

O Correlation ID geralmente é obtido a partir do cabeçalho da requisição `x-correlation-id`. Se este cabeçalho não estiver presente na requisição, o sistema automaticamente gera um novo Correlation ID usando um UUID (Universally Unique Identifier).

Claro, vou ajustar a seção sobre os níveis de log para indicar que o `Fields` é opcional:

---

#### Níveis de Log

Existem diferentes níveis de log que podem ser usados, dependendo da natureza e da severidade da mensagem. Para cada nível de log, você pode incluir uma mensagem e campos adicionais (opcionalmente) para fornecer contexto adicional:

- **Fatal**: Usado para registrar mensagens que indicam um erro crítico que pode resultar na parada da aplicação.

  ```go
  log.Fatal("Mensagem crítica")
  log.Fatal("Mensagem crítica com campos", log.Fields{"chave": "valor"})
  ```

- **Error**: Usado para registrar mensagens de erro que não são críticas, mas que indicam problemas no sistema.

  ```go
  log.Error("Mensagem de erro")
  log.Error("Mensagem de erro com campos", log.Fields{"chave": "valor"})
  ```

- **Warning**: Indica uma situação potencialmente problemática que não é necessariamente um erro.

  ```go
  log.Warning("Mensagem de aviso")
  log.Warning("Mensagem de aviso com campos", log.Fields{"chave": "valor"})
  ```

- **Info**: Usado para mensagens informativas que não representam erros.

  ```go
  log.Info("Mensagem informativa")
  log.Info("Mensagem informativa com campos", log.Fields{"chave": "valor"})
  ```

- **Debug**: Registra detalhes de operações para fins de depuração. Essas mensagens geralmente contêm informações que são úteis durante o desenvolvimento e depuração.

  ```go
  log.Debug("Mensagem de depuração")
  log.Debug("Mensagem de depuração com campos", log.Fields{"chave": "valor"})
  ```

#### Formatação e Campos

O logger suporta a inclusão de campos adicionais nas mensagens de log. Isso é útil para adicionar metadados contextuais às mensagens.

Os campos adicionais podem ser passados usando o tipo `Fields`, que é um mapa de strings para interfaces.

```go
fields := log.Fields{
	"chave1": "valor1",
	"chave2": "valor2",
}
log.Info("Mensagem com campos", fields)
```

#### Ambientes

O logger adapta sua configuração com base nas variáveis de ambiente. Por exemplo, em ambientes de produção, os logs são formatados em JSON, enquanto em ambientes de desenvolvimento são formatados em texto.

