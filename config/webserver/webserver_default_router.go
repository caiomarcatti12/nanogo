package webserver

import (
	"github.com/caiomarcatti12/nanogo/v2/controller"
)

func WebserverDefaultRouter() {
	AddRouter("GET", "/healthcheck", controller.HealthcheckHandler)
}
