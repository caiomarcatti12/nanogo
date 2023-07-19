package main

import (
	"github.com/codelesshub/nanogo/config/env"
	"github.com/codelesshub/nanogo/config/webserver"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	server := webserver.NewWebServer()
	//webserver.AddRouter("GET", "/test", controller.HealthcheckHandler)
	server.Start()
}
