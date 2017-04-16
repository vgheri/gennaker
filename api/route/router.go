package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vgheri/gennaker/api/handler"
)

// NewRouter returns a router with all routes registered
func NewRouter(handler *handler.Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	routes := setupRoutes(handler)
	for _, route := range routes {
		var routeHandler http.Handler
		routeHandler = route.HandlerFunc
		routeHandler = logCall(routeHandler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(routeHandler)
	}

	return router
}

// func accessControl(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
//
// 		if r.Method == "OPTIONS" {
// 			return
// 		}
//
// 		h.ServeHTTP(w, r)
// 	})
// }
