package route

import (
	"net/http"

	"github.com/vgheri/gennaker/api/handler"
)

// Route maps key information for an HTTP route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a collection of routes
type Routes []*Route

func setupRoutes(handler *handler.Handler) Routes {
	return Routes{
		&Route{
			Name:        "CreateDeployment",
			Method:      "POST",
			Pattern:     "/api/v1/deployment",
			HandlerFunc: handler.CreateDeploymentHandler,
		},
		&Route{
			Name:        "NewReleaseNotification",
			Method:      "POST",
			Pattern:     "/api/v1/deployment/newrelease",
			HandlerFunc: handler.NewDeploymentReleaseNotificationHandler,
		},
	}
}
