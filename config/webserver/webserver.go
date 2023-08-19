package webserver

import (
	"encoding/json"
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
		ws.router.Use(PayloadMiddleware)
	})

	return ws
}

func NewWebServer() *WebServer {
	ws := getWebServerInstance()

	WebserverDefaultRouter()

	return ws
}

//func AddRouter(method string, path string, f func(http.ResponseWriter, *http.Request)) {
//
//	getWebServerInstance().router.HandleFunc(path, f).Methods(method)
//}

func AddRouter(method string, path string, f func(ctx *HandlerContext) *APIResponse) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		contextPayload := r.Context().Value("payload")

		// Se você esperar um tipo específico para o payload, faça o type assertion aqui
		payload, _ := contextPayload.(map[string]interface{})

		response := f(&HandlerContext{
			Payload:  payload,
			Headers:  r.Header,
			Request:  r,
			Response: w,
		})

		// Set headers
		for key, value := range response.Headers {
			w.Header().Set(key, value)
		}
		// Send the response
		w.WriteHeader(response.StatusCode)
		if response.Data != nil {
			json.NewEncoder(w).Encode(response.Data)
		}
	}
	getWebServerInstance().router.HandleFunc(path, handlerFunc).Methods(method)
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
