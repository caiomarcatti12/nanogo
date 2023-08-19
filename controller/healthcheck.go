package controller

import (
	"net/http"

	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
	"github.com/caiomarcatti12/nanogo/v2/config/log"
)

func HealthcheckHandler(ctx *webserver.HandlerContext) *webserver.APIResponse {
	logger := log.GetLoggerFromContext(ctx.Request.Context())
	logger.Debug("Healthcheck request received")

	return &webserver.APIResponse{
		Data:       "Service is up and running", // ou simplesmente nil se você não quiser enviar uma mensagem
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain"},
	}
}
