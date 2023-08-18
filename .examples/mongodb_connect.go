package main

import (
	mongodb "github.com/caiomarcatti12/nanogo/config/database"
	"github.com/caiomarcatti12/nanogo/config/env"
)

func main() {
	// Carrega o arquivo .env
	env.LoadEnv()

	// Connecta no mongodb
	mongodb.ConnectMongoDB()
}
