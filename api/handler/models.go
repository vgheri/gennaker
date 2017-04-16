package handler

type CreateDeploymentRequest struct {
	Name          string `json:"name"`
	ChartName     string `json:"chart_name"`
	RepositoryURL string `json:"repository_url"`
}

type CreateDeploymentResponse struct {
	ID int `json:"id"`
}
