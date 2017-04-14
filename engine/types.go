package engine

import "time"

//Deployment is the unit of work of gennaker.
//A deployment maps to a helm chart and contains the history of
//releases to a cluster for this chart
type Deployment struct {
	ID            int        `json:"id"`
	ChartName     string     `json:"chart_name"`
	RepositoryURL string     `json:"repository_url"`
	Releases      []*Release `json:"releases"`
	Pipeline      *Pipeline  `json:"pipeline"`
	CreationDate  time.Time  `json:"creation_date"`
	LastUpdate    time.Time  `json:"last_update"`
}

//Release models a versioned release of the content of an helm chart
type Release struct {
	ID           int       `json:"id"`
	DeploymentID int       `json:"deployment_id"`
	Date         time.Time `json:"date"`
	Namespace    string    `json:"namespace"`
	Image        string    `json:"image"`
	ImageTag     string    `json:"image_tag"`
}

//PipelineStep models a specific step in the deployment lifecycle
type PipelineStep struct {
	ID              int             `json:"id"`
	StepNumber      int             `json:"step_number"`
	ParentID        int             `json:"parent_id"`
	PipelineID      int             `json:"pipeline_id"`
	TargetNamespace string          `json:"target_namespace"`
	AutomaticDeploy bool            `json:"automatic_deploy"`
	NextSteps       []*PipelineStep `json:"next_steps"`
}

//Pipeline models the lifecycle of a deployment
//Ex: new pushes must be automatically deployed to dev and load namespaces,
//then a manual promotion can happen from dev to staging and from staging to prod
type Pipeline struct {
	ID         int             `json:"id"`
	LastUpdate time.Time       `json:"last_update"`
	RootSteps  []*PipelineStep `json:"root_steps"`
}

//DeploymentService describes all functionalities exposed by gennaker
type DeploymentService interface {
	ListDeployments(limit, offset int) ([]*Deployment, error)
	ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error)
	GetDeployment(id int) (*Deployment, error)
}

//DeploymentRepository contains all necessary database support methods
//needed by the DeploymentService
type DeploymentRepository interface {
	ListDeployments(limit, offset int) ([]*Deployment, error)
	ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error)
	GetDeployment(id int) (*Deployment, error)
}
