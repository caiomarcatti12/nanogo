# Web Server

API HTTP construída sobre `gorilla/mux` com middlewares e registro automático de rotas.

## Principais APIs
- `Factory(env, logger, i18n, di, telemetry, contextManager) IWebServer`.
- Métodos: `AddRoute`, `AddMidleware`, `Start`.

## Variáveis de Ambiente
- `WEB_SERVER_HOST`
- `WEB_SERVER_PORT`
- `WEB_SERVER_CERTIFICATE`
- `WEB_SERVER_KEY`
- `WEBSERVER_ACCESS_LOG`
- `WEBSERVER_ORIGINS`, `WEBSERVER_HEADERS`, `WEBSERVER_METHODS`
- `WEB_SERVER_MAX_UPLOAD_SIZE`

## Exemplo de Uso
```go
ws := webserver.Factory(envAdapter, logger, i18n, container, telemetry, cm)
ws.Start()
```
