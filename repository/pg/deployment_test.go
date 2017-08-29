package pg

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/vgheri/gennaker/engine"
)

var pg *pgRepository
var db *sql.DB

func TestMain(m *testing.M) {
	var host, port, user, password, dbname string
	if host = os.Getenv("POSTGRES_HOST"); host == "" {
		host = "127.0.0.1"
	}
	if port = os.Getenv("POSTGRES_PORT"); port == "" {
		port = "5432"
	}
	if user = os.Getenv("POSTGRES_USER"); user == "" {
		user = "postgres"
	}
	if password = os.Getenv("POSTGRES_PASSWORD"); password == "" {
		password = "password"
	}
	if dbname = os.Getenv("POSTGRES_DB"); dbname == "" {
		dbname = "gennaker"
	}

	client, err := NewClient(host, port, user, password, dbname, 250)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pg = client.(*pgRepository)
	if err != nil {
		panic(err)
	}
	db, _ = pg.getDB()
	r := m.Run()
	teardown(db)
	os.Exit(r)
}

func Test_ListDeployments(t *testing.T) {
	teardown(db)
	deployments, err := pg.ListDeployments(10, 0)
	if err != nil {
		t.Fatalf("Error getting deployments. Error details: %v", err)
	}
	if len(deployments) != 0 {
		t.Fatalf("Expected to have 0 deployments with an empty db")
	}
	insertDummyData(db)
	deployments, err = pg.ListDeployments(10, 0)
	if err != nil {
		t.Fatalf("Error getting deployments. Error details: %v", err)
	}
	if len(deployments) != 2 {
		t.Fatalf("Expected to have 2 deployments with an empty db")
	}
}

func Test_GetDeployment(t *testing.T) {

	teardown(db)
	deployment, err := pg.GetDeployment("abc")
	if err != engine.ErrResourceNotFound {
		t.Fatalf("Expected resource not found, got %v", deployment)
	}

	insertDummyData(db)

	deployment, err = pg.GetDeployment(firstTestDeploymentName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if deployment == nil {
		t.Fatalf("Expected deployment not to be nil")
	}
	if deployment.ID == 0 {
		t.Fatalf("Expected deployment id to be > 0")
	}
	if deployment.Name != "test app" {
		t.Fatalf("Expeced name to be test app, got %s", deployment.Name)
	}
	if deployment.ChartName != "test-chart" {
		t.Fatalf("Expected chart name to be test-chart, got %s", deployment.ChartName)
	}
	if deployment.RepositoryURL != "https://test.com/helm/charts" {
		t.Fatalf("Expected repository URL to be `https://test.com/helm/charts`, got %s", deployment.RepositoryURL)
	}
	if len(deployment.Pipeline) != 2 {
		t.Fatalf("Expected pipeline length 2, got %d", len(deployment.Pipeline))
	}
	if deployment.Pipeline[0].ID < 1 ||
		deployment.Pipeline[0].StepNumber != 1 ||
		deployment.Pipeline[0].ParentStepNumber != 0 ||
		deployment.Pipeline[0].TargetNamespace != "dev" ||
		deployment.Pipeline[0].AutomaticDeploy != true ||
		len(deployment.Pipeline[0].NextSteps) != 1 ||
		deployment.Pipeline[0].NextSteps[0].ID == 0 ||
		deployment.Pipeline[0].NextSteps[0].StepNumber != 3 ||
		deployment.Pipeline[0].NextSteps[0].ParentStepNumber != 1 ||
		deployment.Pipeline[0].NextSteps[0].TargetNamespace != "ppd" ||
		deployment.Pipeline[0].NextSteps[0].AutomaticDeploy != false ||
		len(deployment.Pipeline[0].NextSteps[0].NextSteps) != 1 ||
		deployment.Pipeline[0].NextSteps[0].NextSteps[0].ID < 1 ||
		deployment.Pipeline[0].NextSteps[0].NextSteps[0].StepNumber != 4 ||
		deployment.Pipeline[0].NextSteps[0].NextSteps[0].ParentStepNumber != 3 ||
		deployment.Pipeline[0].NextSteps[0].NextSteps[0].TargetNamespace != "prod" ||
		deployment.Pipeline[0].NextSteps[0].NextSteps[0].AutomaticDeploy != false ||
		len(deployment.Pipeline[0].NextSteps[0].NextSteps[0].NextSteps) != 0 ||
		deployment.Pipeline[1].ID == 0 ||
		deployment.Pipeline[1].StepNumber != 2 ||
		deployment.Pipeline[1].ParentStepNumber != 0 ||
		deployment.Pipeline[1].TargetNamespace != "int" ||
		deployment.Pipeline[1].AutomaticDeploy != true ||
		len(deployment.Pipeline[1].NextSteps) != 0 ||
		len(deployment.Releases) != 5 ||
		deployment.Releases[0].ID == 0 ||
		deployment.Releases[0].DeploymentID == 0 ||
		deployment.Releases[0].Name != "happy-panda" ||
		deployment.Releases[0].ImageTag != "0.0.1" ||
		deployment.Releases[0].Namespace != "ppd" ||
		deployment.Releases[0].Values != "a=1" ||
		deployment.Releases[0].Chart != "test-chart" ||
		deployment.Releases[0].Status != engine.Deployed ||
		deployment.Releases[4].ID == 0 ||
		deployment.Releases[4].DeploymentID == 0 ||
		deployment.Releases[4].ImageTag != "0.0.1" ||
		deployment.Releases[4].Namespace != "dev" ||
		deployment.Releases[4].Values != "a=1" {
		t.Fatalf("Malformed pipeline %+v", deployment.Releases[0])
	}
}

func Test_CreateDeployment(t *testing.T) {

	teardown(db)
	err := pg.CreateDeployment(nil)
	if err != engine.ErrInvalidDeployment {
		t.Fatalf("Expected ErrInvalidDeployment creating a nil deployment, got %v", err)
	}

	teardown(db)
	deployment := &engine.Deployment{
		Name:          "unit test app",
		ChartName:     "test",
		ChartVersion:  "0.1.0",
		RepositoryURL: "test",
		Pipeline:      nil,
	}
	err = pg.CreateDeployment(deployment)
	if err != engine.ErrInvalidPipeline {
		t.Fatalf("Expected test to fail with invalid pipeline, got %v", err)
	}

	teardown(db)
	deployment = &engine.Deployment{
		Name:          "unit test app",
		ChartName:     "test",
		ChartVersion:  "0.1.0",
		RepositoryURL: "http://test.com/charts",
		Pipeline: []*engine.PipelineStep{
			&engine.PipelineStep{
				StepNumber:       1,
				ParentStepNumber: 0,
				TargetNamespace:  "dev",
				AutomaticDeploy:  true,
				NextSteps: []*engine.PipelineStep{
					&engine.PipelineStep{
						StepNumber:       3,
						ParentStepNumber: 1,
						TargetNamespace:  "ppd",
						AutomaticDeploy:  false,
						NextSteps: []*engine.PipelineStep{
							&engine.PipelineStep{
								StepNumber:       4,
								ParentStepNumber: 3,
								TargetNamespace:  "prod",
								AutomaticDeploy:  false,
							},
						},
					},
				},
			},
			&engine.PipelineStep{
				StepNumber:       2,
				ParentStepNumber: 0,
				TargetNamespace:  "int",
				AutomaticDeploy:  true,
			},
		},
	}
	err = pg.CreateDeployment(deployment)
	if err != nil {
		t.Fatalf("Expected test to succeed, got err %v", err)
	}
	if deployment.ID == 0 {
		t.Fatalf("Expected deployment ID > 0, got %v", deployment.ID)
	}
}
