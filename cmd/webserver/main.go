package main

import (
	"github.com/caiomarcatti12/nanogo/v1/pkg/di"
	"github.com/caiomarcatti12/nanogo/v1/pkg/nanogo"
	"github.com/caiomarcatti12/nanogo/v1/pkg/webserver"
)

func main() {
	nanogo.Bootstrap()

	ws, err := di.Get[webserver.IWebServer]()

	if err != nil {
		panic(err)
	}

	ws.Start()
}
