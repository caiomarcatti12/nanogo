package main

import (
	"github.com/codelesshub/nanogo/config/env"
	"github.com/codelesshub/nanogo/config/log"
)

func main() {
	// Carrega o arquivo .env
	env.LoadRemoteEnv(env.LoadRemoteEnvParams{
		Host:    env.GetEnv("CLOUD_PROPERTIES_HOST", "http://host.docker.internal:8080"),
		Token:   env.GetEnv("CLOUD_PROPERTIES_TOKEN", ""),
		AppName: env.GetEnv("APP_NAME", "api-bin-runner"),
		Env:     env.GetEnv("APP_NAME", "dev"),
	})

	log.Debug(env.GetEnv("SENHA"))
	//
	//server := webserver.NewWebServer()
	//webserver.AddRouter("GET", "/test", controller.HealthcheckHandler)
	//server.Start()

}
