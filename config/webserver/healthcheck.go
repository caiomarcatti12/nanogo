package webserver

import (
	"net/http"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
)

func HealthcheckHandler(ctx *HandlerContext[any]) (interface{}, error) {
	logger := log.GetLoggerFromContext(ctx.Request.Context())
	logger.Debug("Healthcheck request received")

	return &APIResponse{
		Data:       "Service is up and running", // ou simplesmente nil se você não quiser enviar uma mensagem
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain"},
	}, nil
}
