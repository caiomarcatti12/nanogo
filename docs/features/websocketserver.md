# WebSocket Server

Servidor WebSocket integrado ao WebServer para conexões tempo real.

## Principais APIs
- `Factory(env, logger, i18n, ws, di) IWebSocketServer`.
- Métodos: `AddRoute`, `Start` e `HandleConnections`.

## Variáveis de Ambiente
- `WEBSOCKET_SERVER_LOG_INPUT` habilita logs das mensagens recebidas.

## Exemplo de Uso
```go
wss := websocketserver.Factory(envAdapter, logger, i18n, webSrv, container)
wss.Start()
```
