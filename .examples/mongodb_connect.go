package main

import (
	mongodb "github.com/caiomarcatti12/nanogo/v2/config/database"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	// Connecta no mongodb
	mongodb.ConnectMongoDB()
}
