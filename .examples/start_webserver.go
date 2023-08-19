package main

import (
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/webserver"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	server := webserver.NewWebServer()
	//webserver.AddRouter("GET", "/test", controller.HealthcheckHandler)
	
	server.Start()
}
