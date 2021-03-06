package pg

import (
	"database/sql"
	"fmt"
	// Used by database/sql
	_ "github.com/lib/pq"
)

var firstTestDeploymentID int
var secondTestDeploymentID int
var firstTestDeploymentName = "test app"
var secondTestDeploymentName = "test app 2"

func insertDummyData(db *sql.DB) {
	queries := []string{
		//Deployment 1
		"INSERT INTO deployment(name, chart, repository_url) VALUES('" + firstTestDeploymentName + "', 'test-chart', 'https://test.com/helm/charts') RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(1, NULL, (SELECT id FROM deployment where chart = 'test-chart'), 'dev', true) RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(2, NULL, (SELECT id FROM deployment where chart = 'test-chart'), 'int', true) RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(3, 1, (SELECT id FROM deployment where chart = 'test-chart'), 'ppd', false) RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(4, 3, (SELECT id FROM deployment where chart = 'test-chart'), 'prod', false) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-panda', (SELECT id FROM deployment where chart = 'test-chart'), '0.0.1', 'dev', 'a=1', 'test-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-panda', (SELECT id FROM deployment where chart = 'test-chart'), '0.0.1', 'int', 'a=1', 'test-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-panda', (SELECT id FROM deployment where chart = 'test-chart'), '0.0.2', 'dev', 'a=1', 'test-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-panda', (SELECT id FROM deployment where chart = 'test-chart'), '0.0.2', 'int', 'a=1', 'test-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-panda', (SELECT id FROM deployment where chart = 'test-chart'), '0.0.1', 'ppd', 'a=1', 'test-chart', 1) RETURNING id",
		//Deployment 2
		"INSERT INTO deployment(name, chart, repository_url) VALUES('" + secondTestDeploymentName + "', 'test-new-chart', 'https://test.com/helm/charts') RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(1, NULL, (SELECT id FROM deployment where chart = 'test-new-chart'), 'dev', true) RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(2, NULL, (SELECT id FROM deployment where chart = 'test-new-chart'), 'int', true) RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(3, 1, (SELECT id FROM deployment where chart = 'test-new-chart'), 'ppd', false) RETURNING id",
		"INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy) VALUES(4, 3, (SELECT id FROM deployment where chart = 'test-new-chart'), 'prod', false) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-zebra', (SELECT id FROM deployment where chart = 'test-new-chart'), '0.0.1', 'dev', 'a=1', 'test-new-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-zebra', (SELECT id FROM deployment where chart = 'test-new-chart'), '0.0.1', 'int', 'a=1', 'test-new-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-zebra', (SELECT id FROM deployment where chart = 'test-new-chart'), '0.0.2', 'dev', 'a=1', 'test-new-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-zebra', (SELECT id FROM deployment where chart = 'test-new-chart'), '0.0.2', 'int', 'a=1', 'test-new-chart', 1) RETURNING id",
		"INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, status) VALUES('happy-zebra', (SELECT id FROM deployment where chart = 'test-new-chart'), '0.0.1', 'ppd', 'a=1', 'test-new-chart', 1) RETURNING id",
	}
	for i, q := range queries {
		row := db.QueryRow(q)
		var err error
		if i == 0 {
			err = row.Scan(&firstTestDeploymentID)
		} else if i == 10 {
			err = row.Scan(&secondTestDeploymentID)
		} else {
			var useless int
			err = row.Scan(&useless)
		}
		if err != nil {
			panic(err)
		}
	}
}

// Teardown wipes out existing data
func teardown(db *sql.DB) {
	queries := []string{
		`DELETE FROM pipeline_step`,
		`DELETE FROM release`,
		`DELETE FROM deployment`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			fmt.Printf("Teardown: %s\n", err)
		}
	}
}
