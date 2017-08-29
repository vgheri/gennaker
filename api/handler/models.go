package handler

// POST /api/v1/deployment
// CreateDeployment endpoint
type CreateDeploymentRequest struct {
	Name          string `json:"name"`
	ChartName     string `json:"chart_name"`
	ChartVersion  string `json:"chart_version"`
	RepositoryURL string `json:"repository_url"`
}

type CreateDeploymentResponse struct {
	ID int `json:"id"` // TODO: remove as it's useless
}

// POST /api/v1/deployment/release
// NewRelease endpoint
type NewDeploymentReleaseNotificationRequest struct {
	DeploymentName string `json:"deployment_name"`
	ImageTag       string `json:"image_tag"`
	ReleaseValues  string `json:"release_values"` // --set parameters to helm install/upgrade
}

type NewDeploymentReleaseNotificationResponse struct {
	Reports []string `json:"reports"`
}
