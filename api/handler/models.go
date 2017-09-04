package handler

// CreateDeploymentRequest POST /api/v1/deployment
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

// NewDeploymentReleaseNotificationRequest POST /api/v1/deployment/release
// NewRelease endpoint
type NewDeploymentReleaseNotificationRequest struct {
	DeploymentName string `json:"deployment_name"`
	ImageTag       string `json:"image_tag"`
	ReleaseValues  string `json:"release_values"` // --set parameters to helm install/upgrade
}

type NewDeploymentReleaseNotificationResponse struct {
	Reports []string `json:"reports"`
}

// PromoteReleaseRequest POST /api/v1/deployment/{name}/release/promote
type PromoteReleaseRequest struct {
	FromNamespace string `json:"from_namespace"`
	ReleaseValues string `json:"release_values"`
}

type PromoteReleaseResponse struct {
	Reports []string `json:"reports"`
}
