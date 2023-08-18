package main

import (
	"github.com/caiomarcatti12/nanogo/config/env"
	"github.com/caiomarcatti12/nanogo/config/webserver"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	server := webserver.NewWebServer()
	//webserver.AddRouter("GET", "/test", controller.HealthcheckHandler)
	server.Start()
}
