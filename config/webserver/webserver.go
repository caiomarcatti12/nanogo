package webserver

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/caiomarcatti12/nanogo/config/log"
	"github.com/gorilla/mux"
)

type WebServer struct {
	router *mux.Router
	port   string
}

var (
	once sync.Once
	ws   *WebServer
)

func getWebServerInstance() *WebServer {
	once.Do(func() {
		router := mux.NewRouter()
		port := getPortWebServer()

		ws = &WebServer{
			router: router,
			port:   port,
		}

		ws.router.Use(CorrelationIDMiddleware)
	})

	return ws
}

func NewWebServer() *WebServer {
	ws := getWebServerInstance()

	WebserverDefaultRouter()

	return ws
}

func AddRouter(method string, path string, f func(http.ResponseWriter, *http.Request)) {
	getWebServerInstance().router.HandleFunc(path, f).Methods(method)
}

func (ws *WebServer) Start() {
	port := getPortWebServer()

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
