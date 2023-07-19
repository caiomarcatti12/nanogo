package webserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codelesshub/nanogo/config/log"
	"github.com/gorilla/mux"
)

type WebServer struct {
	router *mux.Router
}

func NewWebServer() *WebServer {
	return &WebServer{
		router: mux.NewRouter(),
	}
}

func (ws *WebServer) AddRouter(path string, router *mux.Router) {
	ws.router.PathPrefix(path).Handler(router)
}

func (ws *WebServer) Start() {
	port := getPortWebServer()

	ws.AddRouter("/", WebServerRouter())

	fmt.Printf("Servidor iniciado em localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, ws.router))
}

func getPortWebServer() string {
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		log.Fatal("A porta do servidor não foi definida no arquivo .env")
	}

	return port
}
