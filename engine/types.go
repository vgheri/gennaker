package engine

import (
	"time"

	"github.com/pkg/errors"
)

//Deployment is the unit of work of gennaker.
//A deployment maps to a helm chart and contains the history of
//releases to a cluster for this chart
//Pipeline models the lifecycle of a deployment
//Ex: new pushes must be automatically deployed to dev and load namespaces,
//then a manual promotion can happen from dev to staging and from staging to prod
type Deployment struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	ChartName     string          `json:"chart_name"`
	RepositoryURL string          `json:"repository_url"`
	Releases      []*Release      `json:"releases"`
	Pipeline      []*PipelineStep `json:"pipeline"`
	CreationDate  time.Time       `json:"creation_date"`
	LastUpdate    time.Time       `json:"last_update"`
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
	ID               int             `json:"id"`
	StepNumber       int             `json:"step_number"`
	ParentStepNumber int             `json:"parent_step_number"`
	DeploymentID     int             `json:"deployment_id"`
	TargetNamespace  string          `json:"target_namespace"`
	AutomaticDeploy  bool            `json:"automatic_deploy"`
	NextSteps        []*PipelineStep `json:"next_steps"`
}

//DeploymentService describes all functionalities exposed by gennaker
type DeploymentEngine interface {
	ListDeployments(limit, offset int) ([]*Deployment, error)
	ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error)
	GetDeployment(id int) (*Deployment, error)
	CreateDeployment(deployment *Deployment) (int, error)
}

//DeploymentRepository contains all necessary database support methods
//needed by the DeploymentService
type DeploymentRepository interface {
	ListDeployments(limit, offset int) ([]*Deployment, error)
	ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error)
	GetDeployment(id int) (*Deployment, error)
	CreateDeployment(deployment *Deployment) error
}

func (d *Deployment) valid() error {
	if d.Name == "" || d.Name == " " {
		return errors.New("Deployment name is invalid")
	}
	if d.ChartName == "" || d.ChartName == " " { // TODO replace with regex
		return errors.New("Chart name is invalid")
	}
	if d.RepositoryURL == "" || d.RepositoryURL == " " { // TODO replace with regex
		return errors.New("Chart repository URL is invalid")
	}
	if d.ID != 0 { // It's an existing deployment
		if d.Pipeline == nil || len(d.Pipeline) == 0 {
			return errors.New("Invalid pipeline")
		}
	}
	return nil
}
