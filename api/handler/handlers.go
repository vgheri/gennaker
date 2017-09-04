package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vgheri/gennaker/engine"
)

// Handler is a strictly typed object containing the list of available handlers
type Handler struct {
	deploymentEngine engine.DeploymentEngine
}

const mimeTypeJSON string = "application/json; charset=UTF-8"

// New initializes the package with the underlying data store instance
func New(e engine.DeploymentEngine) *Handler {
	return &Handler{
		deploymentEngine: e,
	}
}

// APIError models the error object sent back to the client on error
type APIError struct {
	Code    int    //`json:"code"`
	Message string //`json:"message"`
}

// CreateDeploymentHandler creates a deployment
func (h *Handler) CreateDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request
	var reqBody CreateDeploymentRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		writeJSONError(w, err.Error(), 422)
		return
	}
	// Prepare business call
	id, err := h.deploymentEngine.CreateDeployment(
		&engine.Deployment{
			Name:          reqBody.Name,
			ChartName:     reqBody.ChartName,
			ChartVersion:  reqBody.ChartVersion,
			RepositoryURL: reqBody.RepositoryURL,
		})
	if err != nil {
		// TODO: Get the status code from map of errors
		writeJSONError(w, err.Error(),
			http.StatusBadRequest)
		return
	}

	// Encode response
	respBody := CreateDeploymentResponse{ID: id}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		writeJSONError(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// NewDeploymentReleaseNotificationHandler manages the workflow triggered by
// the notification of a new release for a registered deployment
func (h *Handler) NewDeploymentReleaseNotificationHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request
	var reqBody NewDeploymentReleaseNotificationRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		writeJSONError(w, err.Error(), 422)
		return
	}
	// Prepare business call
	reports, err := h.deploymentEngine.HandleNewReleaseNotification(
		&engine.ReleaseNotification{
			DeploymentName: reqBody.DeploymentName,
			ImageTag:       reqBody.ImageTag,
			ReleaseValues:  reqBody.ReleaseValues,
		})
	if err != nil {
		// TODO: Get the status code from map of errors
		writeJSONError(w, err.Error(),
			http.StatusBadRequest)
		return
	}

	// Encode response
	respBody := NewDeploymentReleaseNotificationResponse{Reports: reports}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		writeJSONError(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// PromoteReleaseHandler manages the workflow triggered by
// the promote request
func (h *Handler) PromoteReleaseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deploymentName := vars["name"]
	// Decode request
	var reqBody PromoteReleaseRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		writeJSONError(w, err.Error(), 422)
		return
	}

	// Prepare business call
	reports, err := h.deploymentEngine.PromoteRelease(
		&engine.PromoteRequest{
			DeploymentName: deploymentName,
			FromNamespace:  reqBody.FromNamespace,
			ReleaseValues:  reqBody.ReleaseValues,
		})
	if err != nil {
		// TODO: Get the status code from map of errors
		writeJSONError(w, err.Error(),
			http.StatusBadRequest)
		return
	}

	// Encode response
	respBody := PromoteReleaseResponse{Reports: reports}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		writeJSONError(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func writeJSONError(w http.ResponseWriter, errorMsg string, httpErrorCode int) {
	w.Header().Set("Content-Type", mimeTypeJSON)
	w.WriteHeader(httpErrorCode)
	apiErr := APIError{Message: errorMsg}
	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		// TODO log
	}
}
