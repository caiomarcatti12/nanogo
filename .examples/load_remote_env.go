package main

import (
	"github.com/codelesshub/nanogo/config/env"
	"github.com/codelesshub/nanogo/config/log"
)

func main() {
	//NECESSÁRIO DEFINIR AS VARIAVEIS
	// CLOUD_PROPERTIES_HOST=http://host.docker.internal:8080;
	// APP_NAME=teste;
	// ENV=prod
	// CLOUD_PROPERTIES_TOKEN
	env.LoadRemoteEnv()

	log.Debug(env.GetEnv("SENHA"))
}
