package pg

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	// Used by database/sql
	_ "github.com/lib/pq"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a random string.
func generateRandomString(s int) string {
	b, _ := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)
}

func insertDummyData(db *sql.DB) {
	queries := []string{
		//Deployment 1
		"INSERT INTO pipeline(id) VALUES(1)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(1, 1, NULL, 1, 'dev', true)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(2, 2, NULL, 1, 'int', true)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(3, 3, 1, 1, 'ppd', false)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(4, 4, 3, 1, 'prod', false)",
		"INSERT INTO deployment(id, chart, repository_url, pipeline_id) VALUES(1, 'test-chart', 'https://test.com/helm/charts', 1)",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(1, 1, 'dev', 'testcorp/testimg', 'v0.0.1')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(2, 1, 'int', 'testcorp/testimg', 'v0.0.1')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(3, 1, 'dev', 'testcorp/testimg', 'v0.0.2')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(4, 1, 'int', 'testcorp/testimg', 'v0.0.2')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(5, 1, 'ppd', 'testcorp/testimg', 'v0.0.1')",
		//Deployment 2
		"INSERT INTO pipeline(id) VALUES(2)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(10, 1, NULL, 2, 'dev', true)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(11, 2, NULL, 2, 'int', true)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(12, 3, 1, 2, 'ppd', false)",
		"INSERT INTO pipeline_step(id, step_number, parent_id, pipeline_id, target_namespace, auto_deploy) VALUES(13, 4, 3, 2, 'prod', false)",
		"INSERT INTO deployment(id, chart, repository_url, pipeline_id) VALUES(10, 'test-new-chart', 'https://test.com/helm/charts', 2)",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(10, 10, 'dev', 'testcorp/testimg', 'v0.0.1')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(11, 10, 'int', 'testcorp/testimg', 'v0.0.1')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(12, 10, 'dev', 'testcorp/testimg', 'v0.0.2')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(13, 10, 'int', 'testcorp/testimg', 'v0.0.2')",
		"INSERT INTO release(id, deployment_id, namespace, image, tag) VALUES(14, 10, 'ppd', 'testcorp/testimg', 'v0.0.1')",
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			fmt.Printf("InsertDummyData: %s\n", err)
		}
	}
}

// Teardown wipes out existing data
func teardown(db *sql.DB) {
	queries := []string{
		`DELETE FROM pipeline_step`,
		`DELETE FROM release`,
		`DELETE FROM deployment`,
		`DELETE FROM pipeline`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			fmt.Printf("Teardown: %s\n", err)
		}
	}
}
