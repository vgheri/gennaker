package engine

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
func (r fakeRepository) GetDeployment(id int) (*Deployment, error) {
	return &Deployment{}, nil
}
func (r fakeRepository) CreateDeployment(deployment *Deployment) error {
	return nil
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
	Convey("Testing CreateDeployment()", t, FailureContinues, func() {
		Convey("With an invalid deployment", func() {
			Convey("Should return error ErrInvalidDeployment", func() {
				invalidDeployment := &Deployment{
					Name:          "",
					ChartName:     "",
					RepositoryURL: "https://fakerepository.com",
				}
				_, err := testEngine.CreateDeployment(invalidDeployment)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldStartWith, "Deployment is invalid")
			})
		})
		Convey("With an invalid repository url", func() {
			Convey("Should return error", func() {
				invalidDeployment := &Deployment{
					Name:          "test app",
					ChartName:     "test",
					RepositoryURL: "https://fakerepository.com",
				}
				_, err := testEngine.CreateDeployment(invalidDeployment)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldStartWith, "AddRepository failed")
			})
		})
		Convey("With a chart that is missing the gennaker.yml", func() {
			Convey("Should return error", func() {
				invalidDeployment := &Deployment{
					Name:          "test app",
					ChartName:     "consul",
					RepositoryURL: "https://kubernetes-charts.storage.googleapis.com",
				}
				_, err := testEngine.CreateDeployment(invalidDeployment)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldStartWith, "Build pipeline failed")
			})
		})
	})
}
