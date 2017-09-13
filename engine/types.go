package engine

import (
	"strings"
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
	ChartVersion  string          `json:"chart_version"`
	RepositoryURL string          `json:"repository_url"`
	Releases      []*Release      `json:"releases"`
	Pipeline      []*PipelineStep `json:"pipeline"`
	CreationDate  time.Time       `json:"creation_date"`
	LastUpdate    time.Time       `json:"last_update"`
}

//Release models a versioned release of the content of an helm chart
type Release struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	DeploymentID int                    `json:"deployment_id"`
	ImageTag     string                 `json:"image_tag"`
	Date         time.Time              `json:"date"`
	Namespace    string                 `json:"namespace"`
	Values       string                 `json:"values"`
	Chart        string                 `json:"chart"`
	ChartVersion string                 `json:"chart_version"`
	Revision     int                    `json:"revision"`
	Status       GennakerReleaseOutcome `json:"status"`
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

type ReleaseNotification struct {
	DeploymentName string
	ImageTag       string
	ReleaseValues  string
}

type PromoteRequest struct {
	DeploymentName string
	FromNamespace  string
	ImageTag       string
	ReleaseValues  string
}

type RollbackRequest struct {
	DeploymentName string
	Namespace      string
	Revision       int
}

//DeploymentService describes all functionalities exposed by gennaker
type DeploymentEngine interface {
	ListDeployments(limit, offset int) ([]*Deployment, error)
	ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error)
	GetDeployment(name string) (*Deployment, error)
	CreateDeployment(deployment *Deployment) (int, error)
	HandleNewReleaseNotification(notification *ReleaseNotification) ([]string, error)
	PromoteRelease(request *PromoteRequest) ([]string, error)
	Rollback(request *RollbackRequest) (string, error)
}

//DeploymentRepository contains all necessary database support methods
//needed by the DeploymentService
type DeploymentRepository interface {
	ListDeployments(limit, offset int) ([]*Deployment, error)
	ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error)
	GetDeployment(name string) (*Deployment, error)
	CreateDeployment(deployment *Deployment) error
	CreateRelease(release *Release) (int, error)
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

func (r *ReleaseNotification) valid() error {
	if len(strings.TrimSpace(r.DeploymentName)) == 0 {
		return errors.New("Deployment name cannot be empty")
	}
	if len(strings.TrimSpace(r.ImageTag)) == 0 {
		return errors.New("ImageTag cannot be empty")
	}
	if len(strings.TrimSpace(r.ReleaseValues)) != 0 && // TODO use regex
		!strings.Contains(r.ReleaseValues, "=") {
		return errors.New("Invalid ReleaseValues")
	}
	return nil
}

func (r *PromoteRequest) valid() error {
	if len(strings.TrimSpace(r.DeploymentName)) == 0 {
		return errors.New("Deployment name cannot be empty")
	}
	if len(strings.TrimSpace(r.FromNamespace)) == 0 {
		return errors.New("FromNamespace cannot be empty")
	}
	if len(strings.TrimSpace(r.ReleaseValues)) != 0 && // TODO use regex
		!strings.Contains(r.ReleaseValues, "=") {
		return errors.New("Invalid ReleaseValues")
	}
	return nil
}

func (r *RollbackRequest) valid() error {
	if len(strings.TrimSpace(r.DeploymentName)) == 0 {
		return errors.New("Deployment name cannot be empty")
	}
	if len(strings.TrimSpace(r.Namespace)) == 0 {
		return errors.New("Namespace cannot be empty")
	}
	return nil
}
