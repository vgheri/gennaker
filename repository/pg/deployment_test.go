package pg

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
	Convey("Testing ListDeployments()", t, FailureContinues, func() {
		Convey("With an empty db", func() {
			teardown(db)
			Convey("There should be no deployments", func() {
				deployments, err := pg.ListDeployments(10, 0)
				So(err, ShouldBeNil)
				So(len(deployments), ShouldEqual, 0)
			})
		})
		Convey("With a populated db", func() {
			insertDummyData(db)
			Convey("There should be 2 deployments", func() {
				deployments, err := pg.ListDeployments(10, 0)
				So(err, ShouldBeNil)
				So(len(deployments), ShouldEqual, 2)
			})
		})
	})
}

func Test_GetDeployment(t *testing.T) {
	Convey("Testing GetDeployment(1)", t, FailureContinues, func() {
		Convey("With an empty db", func() {
			teardown(db)
			Convey("It should return ErrResourceNotFound", func() {
				deployment, err := pg.GetDeployment(1)
				So(err, ShouldEqual, engine.ErrResourceNotFound)
				So(deployment, ShouldBeNil)
			})
		})
		Convey("With a populated db", func() {
			insertDummyData(db)
			Convey("It should return a populated object", func() {
				deployment, err := pg.GetDeployment(1)
				So(err, ShouldBeNil)
				So(deployment.ID, ShouldEqual, 1)
				So(deployment.ChartName, ShouldEqual, "test-chart")
				So(deployment.RepositoryURL, ShouldEqual, "https://test.com/helm/charts")
				So(deployment.Pipeline.ID, ShouldEqual, 1)
				So(len(deployment.Pipeline.RootSteps), ShouldEqual, 2)
				So(deployment.Pipeline.RootSteps[0].ID, ShouldEqual, 1)
				So(deployment.Pipeline.RootSteps[0].StepNumber, ShouldEqual, 1)
				So(deployment.Pipeline.RootSteps[0].ParentID, ShouldEqual, 0)
				So(deployment.Pipeline.RootSteps[0].PipelineID, ShouldEqual, 1)
				So(deployment.Pipeline.RootSteps[0].TargetNamespace, ShouldEqual, "dev")
				So(deployment.Pipeline.RootSteps[0].AutomaticDeploy, ShouldBeTrue)
				So(deployment.Pipeline.RootSteps[0].NextSteps, ShouldHaveLength, 1)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].ID, ShouldEqual, 3)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].StepNumber, ShouldEqual, 3)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].ParentID, ShouldEqual, 1)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].PipelineID, ShouldEqual, 1)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].TargetNamespace, ShouldEqual, "ppd")
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].AutomaticDeploy, ShouldBeFalse)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps, ShouldHaveLength, 1)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps[0].ID, ShouldEqual, 4)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps[0].StepNumber, ShouldEqual, 4)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps[0].ParentID, ShouldEqual, 3)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps[0].PipelineID, ShouldEqual, 1)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps[0].TargetNamespace, ShouldEqual, "prod")
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps[0].AutomaticDeploy, ShouldBeFalse)
				So(deployment.Pipeline.RootSteps[0].NextSteps[0].NextSteps[0].NextSteps, ShouldBeEmpty)
				So(deployment.Pipeline.RootSteps[1].ID, ShouldEqual, 2)
				So(deployment.Pipeline.RootSteps[1].StepNumber, ShouldEqual, 2)
				So(deployment.Pipeline.RootSteps[1].ParentID, ShouldEqual, 0)
				So(deployment.Pipeline.RootSteps[1].PipelineID, ShouldEqual, 1)
				So(deployment.Pipeline.RootSteps[1].TargetNamespace, ShouldEqual, "int")
				So(deployment.Pipeline.RootSteps[1].AutomaticDeploy, ShouldBeTrue)
				So(len(deployment.Pipeline.RootSteps[1].NextSteps), ShouldEqual, 0)
				So(deployment.Releases, ShouldHaveLength, 5)
				So(deployment.Releases[0].ID, ShouldEqual, 1)
				So(deployment.Releases[0].DeploymentID, ShouldEqual, 1)
				So(deployment.Releases[0].Namespace, ShouldEqual, "dev")
				So(deployment.Releases[0].Image, ShouldEqual, "testcorp/testimg")
				So(deployment.Releases[0].ImageTag, ShouldEqual, "v0.0.1")
				So(deployment.Releases[4].ID, ShouldEqual, 5)
				So(deployment.Releases[4].DeploymentID, ShouldEqual, 1)
				So(deployment.Releases[4].Namespace, ShouldEqual, "ppd")
				So(deployment.Releases[4].Image, ShouldEqual, "testcorp/testimg")
				So(deployment.Releases[4].ImageTag, ShouldEqual, "v0.0.1")
			})
		})
	})
}
