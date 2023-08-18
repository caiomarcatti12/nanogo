package webserver

import (
	"github.com/caiomarcatti12/nanogo/controller"
)

func WebserverDefaultRouter() {
	AddRouter("GET", "/healthcheck", controller.HealthcheckHandler)
}
