package infrastructurec

import (
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methos", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func SetRoutes(ec *Infrastructure.EventController) {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", MethodHandler)
	http.Handle("/", corsMiddleware(mux))

}
