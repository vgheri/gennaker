package helm

import (
	"bufio"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/vgheri/gennaker/utils"
)

const helmCmd = "helm"

// ReleaseStatus models different statutes used by helm
// to report the outcome of an operation that manages a release
type ReleaseStatus string

func convertToHelmReleasStatus(status string) ReleaseStatus {
	var releaseStatus ReleaseStatus
	switch status {
	case "DEPLOYED":
		releaseStatus = Deployed
	case "DELETED":
		releaseStatus = Deleted
	case "SUPERSEDED":
		releaseStatus = Superseded
	case "FAILED":
		releaseStatus = Failed
	default:
		releaseStatus = Unknown
	}
	return releaseStatus
}

// Helm release statutes
const (
	Unknown    ReleaseStatus = "UNKNOWN"
	Deployed                 = "DEPLOYED"
	Deleted                  = "DELETED"
	Superseded               = "SUPERSEDED"
	Failed                   = "FAILED"
	Deleting                 = "DELETING"
)

// GetRepositoryName retrieves the name of an installed repository by URL.
// Returns an empty string if the repository is not installed
// TODO should probably return error in case no repository is found
func GetRepositoryName(repositoryURL string) (string, error) {
	// helm repo list
	cmdName := helmCmd
	cmdArgs := []string{"repo", "list"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.Wrap(err, "Error creating StdoutPipe for helm repo list")
		return "", err
	}
	var repositoryName string
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, repositoryURL) { // TODO improve with regex
				i := strings.Index(line, "http")
				repositoryName = strings.TrimRight(line[:i-1], " ")
				break
			}
		}
	}()
	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "Could not start command helm repo list:")
	}

	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, "Error waiting for Cmd")
	}
	return repositoryName, nil
}

// Fetch attempts to download and unpack the remote chart into the desired location.
// If version is not provided, than latest version will be downloaded.
// Returns the path to the chart or an error
func Fetch(repositoryName, chartName, version, savePath string) (string, error) {
	if repositoryName == "" || chartName == "" {
		return "", errors.New("Failed at fetching chart: repository name and chart name are mandatory.")
	}
	cmdName := helmCmd
	var cmdArgs []string
	pkg := fmt.Sprintf("%s/%s", repositoryName, chartName)
	if version == "" {
		cmdArgs = []string{"fetch", "-d", savePath, "--untar", pkg}
	} else {
		cmdArgs = []string{"fetch", "-d", savePath, "--untar", "--version", version, pkg}
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		errors.Wrap(err, "Error creating StdoutPipe for helm fetch")
		return "", err
	}
	scanner := bufio.NewScanner(cmdReader)
	fetchSuccess := true
	var helmErrorMsg string
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Error:") {
				fetchSuccess = false
				helmErrorMsg = line
				break
			}
		}
	}()
	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "Could not start command helm fetch:")
	}

	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, "Error waiting for command helm fetch")
	}
	if !fetchSuccess {
		return "", errors.Errorf("Failed at fetching chart: %s", helmErrorMsg)
	}
	return path.Join(savePath, chartName), nil
}

// AddRepository attemps to add a helm repository.
// The name is randomly generated.
func AddRepository(url string) (string, error) {
	// helm repo add
	name := generateRandomRepoName()

	cmdName := helmCmd
	cmdArgs := []string{"repo", "add", name, url}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.Wrap(err, "Error creating StdoutPipe for helm repo add")
		return "", err
	}
	scanner := bufio.NewScanner(cmdReader)
	addSuccess := true
	var helmErrorMsg string
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Error:") {
				addSuccess = false
				helmErrorMsg = line
				break
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "Could not start command helm repo add")
	}
	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, "Error waiting for command helm repo add")
	}
	if !addSuccess {
		return "", errors.Errorf("Failed at adding repository with URL %s. Details: %s", url, helmErrorMsg)
	}
	return name, nil
}

// InstallOrUpgrade installs or upgrades a given release name for the specified chart into the desired namespace.
// If no prior release with the given releaseName is found, an install will be performed, an upgrade otherwise.
func InstallOrUpgrade(releaseName, namespace, repositoryName, chartName, valuesFilePath, releaseValues string) (string, error) {
	if len(strings.TrimSpace(repositoryName)) == 0 || len(chartName) == 0 {
		return "", errors.New("Repository name and chart name are mandatory")
	}
	if len(strings.TrimSpace(releaseName)) == 0 {
		return "", errors.New("Release name is mandatory")
	}
	cmdName := helmCmd
	var cmdArgs = []string{"upgrade", "-i"}
	if len(strings.TrimSpace(namespace)) != 0 {
		cmdArgs = append(cmdArgs, "--namespace", namespace)
	}
	if len(strings.TrimSpace(valuesFilePath)) != 0 {
		cmdArgs = append(cmdArgs, "-f", valuesFilePath)
	}
	if len(strings.TrimSpace(releaseValues)) != 0 {
		cmdArgs = append(cmdArgs, "--set", releaseValues)
	}
	cmdArgs = append(cmdArgs, releaseName)
	pkg := fmt.Sprintf("%s/%s", repositoryName, chartName)
	cmdArgs = append(cmdArgs, pkg)

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		errors.Wrap(err, "Error creating StdoutPipe for helm upgrade")
		return "", err
	}
	scanner := bufio.NewScanner(cmdReader)
	var output string
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			output = strings.Join([]string{output, line}, "\n")
		}
	}()
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return "", errors.Wrap(err, "Error creating StderrPipe for helm install")
	}
	errScanner := bufio.NewScanner(cmdErrReader)
	installSuccess := true
	var helmErrorMsg string
	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			if strings.HasPrefix(line, "ERROR:") {
				installSuccess = false
				helmErrorMsg = line
				break
			}
		}
	}()
	fmt.Printf("%s %s\n", cmdName, cmdArgs)
	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "Could not start command helm upgrade:")
	}
	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Error waiting for command helm upgrade: %s", helmErrorMsg))
	}
	if !installSuccess {
		return output, errors.Errorf("Failed at installing chart: %s", helmErrorMsg)
	}
	return output, nil
}

// Status wraps the helm status command.
// Returns the status of the release using one of possible helm statuses,
// the output of the command and the error, if any
func Status(releaseName string) (ReleaseStatus, string, error) {
	if releaseName == "" {
		return Unknown, "", errors.New("Failed at fetching release status: release name is mandatory")
	}
	cmdName := helmCmd
	cmdArgs := []string{"status", releaseName}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		errors.Wrap(err, "Error creating StdoutPipe for helm status")
		return Unknown, "", err
	}
	scanner := bufio.NewScanner(cmdReader)
	var releaseStatus ReleaseStatus
	var getStatusFailed bool
	var output string
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			output = fmt.Sprintf("%s\n%s", output, line)
			if strings.HasPrefix(line, "STATUS:") {
				parts := strings.Split(line, ": ")
				if len(parts) != 2 {
					getStatusFailed = true
					break
				}
				releaseStatus = convertToHelmReleasStatus(parts[1])
			}
		}
	}()
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return Unknown, output, errors.Wrap(err, "Error creating StderrPipe for helm install")
	}
	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			output = fmt.Sprintf("%s\n%s", output, line)
			if strings.HasPrefix(line, "ERROR:") {
				getStatusFailed = true
				break
			}
		}
	}()
	err = cmd.Start()
	if err != nil {
		return Unknown, output, errors.Wrap(err, "Could not start command helm status")
	}
	err = cmd.Wait()
	if err != nil {
		return Unknown, output, errors.Wrap(err, "Error waiting for command helm status")
	}

	if getStatusFailed {
		return Unknown, output, errors.Errorf("Failed at fetching status for release %s", releaseName)
	}

	return releaseStatus, output, nil
}

func generateRandomRepoName() string {
	return utils.GenerateRandomString(6)
}
