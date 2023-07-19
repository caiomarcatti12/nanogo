package webserver

import (
	"github.com/codelesshub/nanogo/controller"
)

func WebserverDefaultRouter() {
	AddRouter("GET", "/healthcheck", controller.HealthcheckHandler)
}
