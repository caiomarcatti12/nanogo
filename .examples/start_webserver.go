package main

import (
	"github.com/codelesshub/nanogo/config/env"
	"github.com/codelesshub/nanogo/config/webserver"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	server := webserver.NewWebServer()

	//caso desejar utilizar outras rotas
	//runRouter := config.WebServerRouter()
	//server.AddRouter("/run", runRouter)
	//
	//otherRouter := outras_rotas.OutrasRotasRouter()
	//server.AddRouter("/", otherRouter)

	server.Start()
}
