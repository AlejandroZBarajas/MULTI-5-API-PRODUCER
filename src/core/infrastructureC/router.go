package infrastructureC

import (
	"minimulti/src/events/infrastructure"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MethodHandler(w http.ResponseWriter, r *http.Request, ec *infrastructure.EventController) {
	switch r.Method {

	case http.MethodPost:
		ec.CreateNewHandler(w, r)

	case http.MethodGet:
		ec.GetAllHandler(w, r)

	case http.MethodDelete:
		ec.DeleteAllHandler(w, r)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func SetRoutes(ec *infrastructure.EventController) {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		MethodHandler(w, r, ec)
	})
	http.Handle("/", corsMiddleware(mux))

}
