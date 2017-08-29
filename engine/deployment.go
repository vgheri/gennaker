package engine

import (
	"fmt"
	"path"

	"github.com/pkg/errors"
	"github.com/vgheri/gennaker/helm"
)

func (e *engine) CreateDeployment(deployment *Deployment) (int, error) {
	if err := deployment.valid(); err != nil {
		return 0, errors.Wrap(err, "Deployment is invalid")
	}
	// 1. Get repository name or add it if non existent
	repoName, err := helm.GetRepositoryName(deployment.RepositoryURL)
	if err != nil {
		return 0, errors.Wrap(err, "GetRepositoryName failed")
	}
	if repoName == "" {
		repoName, err = helm.AddRepository(deployment.RepositoryURL)
		if err != nil {
			return 0, errors.Wrap(err, "AddRepository failed")
		}
	}
	// 2. Retrieve the chart
	saveDir := path.Join(e.chartsDir, deployment.Name)
	fmt.Printf("Save dir: %s\n", saveDir)
	pathToChart, err := helm.Fetch(repoName, deployment.ChartName, "", saveDir)
	if err != nil {
		return 0, errors.Wrap(err, "Fetch chart failed")
	}
	// 3. Retrieve the gennaker.yml file
	pipeline, err := buildPipeline(pathToChart)
	if err != nil {
		return 0, errors.Wrap(err, "Build pipeline failed")
	}
	// 4. Populate the db
	deployment.Pipeline = pipeline
	err = e.db.CreateDeployment(deployment)
	if err != nil {
		return 0, err
	}
	// 5. Remove the directory with the downloaded chart --> We need the chart so do not remove it
	// err = os.RemoveAll(pathToChart)
	// if err != nil {
	// 	return 0, errors.Wrap(err, "Cannot remove chart folder")
	// }
	return deployment.ID, nil
}

func (e *engine) ListDeployments(limit, offset int) ([]*Deployment, error) {
	return nil, nil
}

func (e *engine) ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error) {
	return nil, nil
}
func (e *engine) GetDeployment(name string) (*Deployment, error) {
	return nil, nil
}
