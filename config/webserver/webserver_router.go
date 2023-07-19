package webserver

import (
	"github.com/codelesshub/nanogo/controller"

	"github.com/gorilla/mux"
)

func WebServerRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controller.HealthcheckHandler).Methods("GET")

	return router
}
