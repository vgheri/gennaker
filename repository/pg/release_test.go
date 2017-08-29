package pg

import (
	"testing"

	"github.com/vgheri/gennaker/engine"
)

func Test_CreateRelease(t *testing.T) {
	teardown(db)
	insertDummyData(db)
	tt := []struct {
		testName  string
		release   engine.Release
		shouldErr bool
	}{
		{testName: "Create should succeed", release: engine.Release{Name: "happy-panda", DeploymentID: firstTestDeploymentID, ImageTag: "0.0.1", Namespace: "prod", Values: "dbname=test",
			Chart: "test-chart", ChartVersion: "0.1.0", Status: engine.Deployed}},
		{testName: "Create should succeed with empty values", release: engine.Release{Name: "happy-panda", DeploymentID: firstTestDeploymentID, ImageTag: "0.0.1", Namespace: "dev",
			Chart: "test-chart", Status: engine.Deployed}},
		{testName: "Create should fail with non existent deployment id", release: engine.Release{Name: "happy-panda", DeploymentID: 25689, ImageTag: "0.0.1", Namespace: "dev",
			Status: engine.Deployed}, shouldErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			id, err := pg.CreateRelease(&tc.release)
			if tc.shouldErr {
				if err == nil {
					t.Fatalf("Expected test to fail")
				}
			} else {
				if id == 0 {
					t.Fatalf("Expected id to be > 0")
				}
			}
		})
	}
}
