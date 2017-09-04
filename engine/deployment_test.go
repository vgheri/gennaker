package engine

import (
	"os"
	"path"
	"strings"
	"testing"
)

type fakeRepository struct{}

var repository fakeRepository
var testEngine DeploymentEngine

func (r fakeRepository) ListDeployments(limit, offset int) ([]*Deployment, error) {
	return []*Deployment{}, nil
}
func (r fakeRepository) ListDeploymentsWithStatus(limit, offset int) ([]*Deployment, error) {
	return []*Deployment{}, nil
}
func (r fakeRepository) GetDeployment(name string) (*Deployment, error) {
	return &Deployment{}, nil
}
func (r fakeRepository) CreateDeployment(deployment *Deployment) error {
	return nil
}
func (r fakeRepository) CreateRelease(release *Release) (int, error) {
	return 0, nil
}

func TestMain(m *testing.M) {
	repository := &fakeRepository{}
	var chartsFolder string
	if chartsFolder = os.Getenv("CHARTS_FOLDER"); chartsFolder == "" {
		gopath := os.Getenv("GOPATH")
		chartsFolder = path.Join(gopath, "src", "github.com", "vgheri", "gennaker", "charts")
	}

	testEngine = New(repository, chartsFolder)
	r := m.Run()
	os.RemoveAll(chartsFolder)
	os.Exit(r)
}

func Test_CreateDeployment(t *testing.T) {

	invalidDeployment := &Deployment{
		Name:          "",
		ChartName:     "",
		RepositoryURL: "https://fakerepository.com",
	}
	_, err := testEngine.CreateDeployment(invalidDeployment)
	if !strings.HasPrefix(err.Error(), "Deployment is invalid") {
		t.Fatalf("Expected invalid deployment, got nothing")
	}

	invalidDeployment = &Deployment{
		Name:          "test app",
		ChartName:     "test",
		RepositoryURL: "https://fakerepository.com",
	}
	_, err = testEngine.CreateDeployment(invalidDeployment)
	if !strings.HasPrefix(err.Error(), "AddRepository failed") {
		t.Fatalf("Expected error AddRepository failed, got %v", err)
	}

	invalidDeployment = &Deployment{
		Name:          "test app",
		ChartName:     "consul",
		RepositoryURL: "https://kubernetes-charts.storage.googleapis.com",
	}
	_, err = testEngine.CreateDeployment(invalidDeployment)
	if !strings.HasPrefix(err.Error(), "Build pipeline failed") {
		t.Fatalf("Expected error Build pipeline failed, got %v", err)
	}
}

func Test_GetDeployment(t *testing.T) {
	invalidDeploymentName := "   "
	_, err := testEngine.GetDeployment(invalidDeploymentName)
	if !strings.HasPrefix(err.Error(), "A non empty deployment name is mandatory") {
		t.Fatalf("Expected invalid deployment name, got nothing")
	}
}
