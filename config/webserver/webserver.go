package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/caiomarcatti12/nanogo/v2/config/errors"
	"net/http"
	"os"
	"sync"

	"github.com/caiomarcatti12/nanogo/v2/config/log"
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

// No pacote webserver

func AddRouter(method string, path string, f func(ctx *HandlerContext) (interface{}, error)) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		contextPayload := r.Context().Value("payload")

		// Se você espera um tipo específico para o payload, faça o type assertion aqui
		payload, _ := contextPayload.(map[string]interface{})

		data, err := f(&HandlerContext{
			Payload:  payload,
			Headers:  r.Header,
			Request:  r,
			Response: w,
		})

		if err != nil {
			if customErr, ok := err.(*errors.CustomError); ok {
				http.Error(w, customErr.Message, customErr.Code)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Verifique se a resposta é uma APIResponse
		if apiResponse, ok := data.(*APIResponse); ok {
			// Se for uma APIResponse, envie diretamente
			for key, value := range apiResponse.Headers {
				w.Header().Set(key, value)
			}
			w.WriteHeader(apiResponse.StatusCode)
			if apiResponse.Data != nil {
				json.NewEncoder(w).Encode(apiResponse.Data)
			}
		} else {
			// Caso contrário, envolva os dados em uma resposta padrão
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(data)
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
