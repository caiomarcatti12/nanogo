package main

import (
	"github.com/caiomarcatti12/nanogo/pkg/di"
	"github.com/caiomarcatti12/nanogo/pkg/nanogo"
	"github.com/caiomarcatti12/nanogo/pkg/webserver"
)

func main() {
	nanogo.Bootstrap()

	ws, err := di.Get[webserver.IWebServer]()

	if err != nil {
		panic(err)
	}

	ws.Start()
}
