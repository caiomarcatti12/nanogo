package webserver

import (
	"context"
	"encoding/json"
	"net/http"
)

func PayloadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var payload map[string]interface{} // Ou o tipo que você espera

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newCtx := context.WithValue(ctx, "payload", payload)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
