package infrastructurec

import "net/http"

/* import (
	"MINI-MULTI-API-PRODUCER/src/events/infrastructure"
) */

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		ec.PostHandler(w, r)

	case http.MethodGet:
		ec.GetAllHandler(w, r)

	case http.MethodDelete:
		ec.DeleteHandler(w, r)
	}
}

func SetRoutes() {

}
