package helm

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// GetRepositoryName retrieves the name of an installed repository by URL.
// Returns an empty string if the repository is not installed
// TODO should probably return error in case no repository is found
func GetRepositoryName(repositoryURL string) (string, error) {
	// helm repo list
	cmdName := "helm"
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
		return "", errors.Wrap(err, "Could not satrt command helm repo list:")
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
	cmdName := "helm"
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
		errors.Wrap(err, "Error creating StdoutPipe for helm repo list")
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
		return "", errors.Wrap(err, "Could not satrt command helm fetch:")
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
	name, err := generateRandomRepoName()
	if err != nil {
		return "", err
	}
	cmdName := "helm"
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
		return "", errors.Wrap(err, "Could not satrt command helm repo add")
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

func generateRandomRepoName() (string, error) {
	var length = 10
	b := make([]byte, length)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
