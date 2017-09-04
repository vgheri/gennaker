package engine

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vgheri/gennaker/helm"
	"github.com/vgheri/gennaker/utils"
)

// GennakerReleaseOutcomees models different statutes used by gennaker
// to report the outcome of an operation that manages a release
type GennakerReleaseOutcome uint8

const imageTag = "ImageTag"

const (
	Unknown  GennakerReleaseOutcome = 0
	Deployed                        = 1
	Failed                          = 2
)

func (e *engine) HandleNewReleaseNotification(notification *ReleaseNotification) ([]string, error) {
	if notification == nil {
		return nil, ErrInvalidReleaseNotification
	}
	if err := notification.valid(); err != nil {
		return nil, err
	}
	var reports []string
	d, err := e.db.GetDeployment(notification.DeploymentName) // TODO: use e.GetDeployment when it's done
	if err != nil {
		return nil, err
	}
	repoName, err := helm.GetRepositoryName(d.RepositoryURL)
	if err != nil {
		return reports, errors.Wrap(err, fmt.Sprintf("Cannot get repository name for url %s", d.RepositoryURL))
	}
	for _, step := range d.Pipeline {
		releaseNameForNamespace := getReleaseName(d, step)

		// Namespace dependent configuration values are stored in $namespace-values.yml
		// inside the chart located in engine.chartsDir
		namespaceValuesFilePath := getNamespaceValuesFilePath(e.chartsDir, d.Name, d.ChartName, step.TargetNamespace)
		releaseValues := buildReleaseValues(notification.ImageTag, notification.ReleaseValues)
		report, err := helm.InstallOrUpgrade(releaseNameForNamespace, step.TargetNamespace,
			repoName, d.ChartName, namespaceValuesFilePath, releaseValues)
		if err != nil {
			return reports, errors.Wrap(err,
				fmt.Sprintf("Failed at installing or upgrading release %s in namespace %s", releaseNameForNamespace, step.TargetNamespace))
		}
		reports = append(reports, report)
		go finalizeRelease(e.db, d, step.TargetNamespace, releaseNameForNamespace,
			notification.ImageTag, notification.ReleaseValues)
	}
	return reports, nil
}

func (e *engine) PromoteRelease(request *PromoteRequest) ([]string, error) {
	if request == nil {
		return nil, ErrInvalidReleaseNotification
	}
	if err := request.valid(); err != nil {
		return nil, errors.Wrap(err, "Promote request is invalid")
	}
	var reports []string
	d, err := e.db.GetDeployment(request.DeploymentName) // TODO: use e.GetDeployment when it's done
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get deployment")
	}
	repoName, err := helm.GetRepositoryName(d.RepositoryURL)
	if err != nil {
		return reports, errors.Wrap(err, fmt.Sprintf("Cannot get repository name for url %s", d.RepositoryURL))
	}
	pipeline := getPipelineForNamespace(request.FromNamespace, d.Pipeline)
	if len(pipeline) == 0 {
		return nil, errors.Errorf("Cannot promote from namespace %s", request.FromNamespace)
	}
	lastRelease := getLastReleaseForNamespace(request.FromNamespace, d)
	if lastRelease == nil {
		return nil, errors.Errorf("Cannot promote: no release found for namespace %s", request.FromNamespace)
	}
	for _, step := range pipeline {
		releaseNameForNamespace := generateReleaseName(d.Name, request.FromNamespace)
		// Namespace dependent configuration values are stored in $namespace-values.yml
		// inside the chart located in engine.chartsDir
		namespaceValuesFilePath := getNamespaceValuesFilePath(e.chartsDir, d.Name, d.ChartName, request.FromNamespace)
		releaseValues := buildReleaseValues(lastRelease.ImageTag, request.ReleaseValues)
		report, err := helm.InstallOrUpgrade(releaseNameForNamespace, step.TargetNamespace,
			repoName, d.ChartName, namespaceValuesFilePath, releaseValues)
		if err != nil {
			return reports, errors.Wrap(err,
				fmt.Sprintf("Failed at installing or upgrading release %s in namespace %s", releaseNameForNamespace, step.TargetNamespace))
		}
		reports = append(reports, report)
		go finalizeRelease(e.db, d, step.TargetNamespace, releaseNameForNamespace,
			lastRelease.ImageTag, request.ReleaseValues)
	}

	return reports, nil
}

func getReleaseName(d *Deployment, step *PipelineStep) string {
	var releaseNameForNamespace string
	// releases are ordered by most recent to less recent
	for _, r := range d.Releases {
		if r.Namespace == step.TargetNamespace {
			releaseNameForNamespace = r.Name
			break
		}
	}
	if len(releaseNameForNamespace) == 0 {
		releaseNameForNamespace = generateReleaseName(d.Name, step.TargetNamespace)
	}
	return releaseNameForNamespace
}

func getNamespaceValuesFilePath(generalchartsDirPath, deploymentName, chartName, namespace string) string {
	chartPath := path.Join(generalchartsDirPath, deploymentName, chartName)
	return path.Join(chartPath, fmt.Sprintf("%s-values.yaml", namespace))
}

func buildReleaseValues(tag, releaseValues string) string {
	imageTagValue := fmt.Sprintf("%s=%s", imageTag, tag)
	if len(strings.TrimSpace(releaseValues)) == 0 {
		releaseValues = imageTagValue
	} else {
		releaseValues = fmt.Sprintf("%s,%s", releaseValues, imageTagValue)
	}
	return releaseValues
}

func generateReleaseName(deploymentName, namespace string) string {
	// return fmt.Sprintf("%s-%s-%s", deploymentName, namespace, utils.GenerateRandomString(5))
	return fmt.Sprintf("%s-%s", utils.GenerateRandomString(5), utils.GenerateRandomString(5))
}

// finalizeRelease loops for 5 minutes waiting to have a status != Unknown
// to persist release status in db
func finalizeRelease(repository DeploymentRepository, deployment *Deployment, namespace, releaseName, imageTag, releaseValues string) {
	// TODO
	// loop for 5 minutes for status to report either success or failure
	// once it's done, update the DB
	start := time.Now()
	var releaseOutcome GennakerReleaseOutcome
	releaseOutcome = Unknown
	for {
		if time.Since(start) > 5*time.Minute {
			releaseOutcome = Unknown
			break
		}
		status, _, err := helm.Status(releaseName)
		// Release in progress
		if err == nil && status == helm.Unknown {
			continue
		}
		switch status {
		case helm.Deleted, helm.Deleting, helm.Superseded, helm.Deployed:
			releaseOutcome = Deployed
		case helm.Failed:
			releaseOutcome = Failed
		default:
			releaseOutcome = Unknown
		}
		if releaseOutcome != Unknown {
			break
		}
		time.Sleep(20 * time.Second)
	}

	release := &Release{
		Name:         releaseName,
		DeploymentID: deployment.ID,
		ImageTag:     imageTag,
		Date:         time.Now(),
		Namespace:    namespace,
		Values:       releaseValues,
		Chart:        deployment.ChartName,
		ChartVersion: deployment.ChartVersion,
		Status:       releaseOutcome,
	}
	// TODO: log error
	_, _ = repository.CreateRelease(release)
}

func getPipelineForNamespace(namespace string, pipeline []*PipelineStep) []*PipelineStep {
	for _, step := range pipeline {
		if step.TargetNamespace == namespace {
			return step.NextSteps
		} else {
			return getPipelineForNamespace(namespace, step.NextSteps)
		}
	}
	return nil
}

func getLastReleaseForNamespace(namespace string, d *Deployment) *Release {
	for _, r := range d.Releases {
		if r.Namespace == namespace {
			return r
		}
	}
	return nil
}
