package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vgheri/gennaker/api/handler"
	"github.com/vgheri/gennaker/api/route"
	"github.com/vgheri/gennaker/engine"
)

type Server struct {
	deploymentEngine engine.DeploymentEngine
}

func New(engine engine.DeploymentEngine) (*Server, error) {
	if engine == nil {
		return nil, errors.New("Nil deployment engine")
	}
	return &Server{deploymentEngine: engine}, nil
}

func (s *Server) Start(port int32) error {
	handlers := handler.New(s.deploymentEngine)
	router := route.NewRouter(handlers)
	http.Handle("/", accessControl(router))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
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
