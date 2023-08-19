package webserver

import (
	"context"
	"encoding/json"
	"net/http"
	"io"
)

func PayloadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Se for uma solicitação GET, capture os parâmetros de consulta.
		if r.Method == http.MethodGet {
			params := r.URL.Query()
			newCtx := context.WithValue(ctx, "params", params)
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		// Se o corpo da requisição estiver vazio, simplesmente chame o próximo manipulador.
		if r.Body == http.NoBody {
			next.ServeHTTP(w, r)
			return
		}

		var payload map[string]interface{} 

		// Se houver um erro ao decodificar e o erro não for devido a um corpo vazio, retorne um erro.
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if payload != nil {
			newCtx := context.WithValue(ctx, "payload", payload)
			next.ServeHTTP(w, r.WithContext(newCtx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
