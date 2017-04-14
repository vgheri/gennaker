package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func serve(port string) error {
	router := mux.NewRouter().StrictSlash(true)
	// router.Handle("/release/list", nil).Methods("GET")
	http.Handle("/", accessControl(router))
	return http.ListenAndServe(":"+port, nil)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
