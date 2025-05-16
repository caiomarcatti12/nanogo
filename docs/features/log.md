# Logger

Sistema de logs baseado no Logrus com campos adicionais e integração ao contexto.

## Principais APIs
- `Factory(env, contextManager) ILog` cria o logger configurado.
- Métodos de `ILog`: `Info`, `Error`, `Debug`, `Trace`, e variações `*f`.

## Variáveis de Ambiente
- `LOG_LEVEL` define o nível de log.
- `ENV` e `APP_NAME` utilizados nos campos padrão.
- `VERSION` para identificar a versão do serviço.

## Exemplo de Uso
```go
logger := log.Factory(envAdapter, contextManager)
logger.Infof("iniciado")
```
