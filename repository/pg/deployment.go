package pg

import (
	"database/sql"
	"time"

	"github.com/vgheri/gennaker/engine"
)

func (r *pgRepository) ListDeployments(limit, offset int) ([]*engine.Deployment, error) {
	query := `SELECT id, pipeline_id, chart, repository_url, creation_date, last_update
  FROM deployment
  ORDER BY chart LIMIT $1 OFFSET $2;`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	deployments := []*engine.Deployment{}
	for rows.Next() {
		var id, pipelineID int
		var chart, repositoryURL string
		var creationDate, lastUpdate time.Time
		err = rows.Scan(&id, &pipelineID, &chart, &repositoryURL, &creationDate, &lastUpdate)
		if err != nil {
			return nil, err
		}
		deployment := &engine.Deployment{
			ID:            id,
			Pipeline:      &engine.Pipeline{ID: id},
			ChartName:     chart,
			RepositoryURL: repositoryURL,
			CreationDate:  creationDate,
			LastUpdate:    lastUpdate,
		}
		deployments = append(deployments, deployment)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return deployments, nil
}

func (r *pgRepository) ListDeploymentsWithStatus(limit, offset int) ([]*engine.Deployment, error) {
	return nil, nil
}

func (r *pgRepository) GetDeployment(id int) (*engine.Deployment, error) {
	query := `SELECT d.pipeline_id, d.chart, d.repository_url, d.creation_date,
  d.last_update, p.last_update
  FROM deployment d
  INNER JOIN pipeline p on p.id = d.pipeline_id
  WHERE d.id = $1`

	row := r.db.QueryRow(query, id)
	var pipelineID int
	var chart, repositoryURL string
	var creationDate, lastUpdate, pipelineLastUpdate time.Time
	err := row.Scan(&pipelineID, &chart, &repositoryURL, &creationDate,
		&lastUpdate, &pipelineLastUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, engine.ErrResourceNotFound
		}
		return nil, err
	}

	query = `SELECT id, timestamp, namespace, image, tag
	FROM release
	WHERE deployment_id = $1
	ORDER BY timestamp asc;`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	releases := []*engine.Release{}
	for rows.Next() {
		var releaseID int
		var timestamp time.Time
		var namespace, image, tag string
		err = rows.Scan(&releaseID, &timestamp, &namespace, &image, &tag)
		if err != nil {
			return nil, err
		}
		release := &engine.Release{
			ID:           releaseID,
			DeploymentID: id,
			Date:         timestamp,
			Namespace:    namespace,
			Image:        image,
			ImageTag:     tag,
		}
		releases = append(releases, release)
	}

	// Build the tree of pipeline steps
	query = `SELECT id, step_number, parent_id, target_namespace, auto_deploy
  FROM pipeline_step
  WHERE pipeline_id = $1
  ORDER BY step_number asc;`
	rows, err = r.db.Query(query, pipelineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stepsMap := make(map[int]*engine.PipelineStep)
	for rows.Next() {
		var stepID, stepNumber, parentID int
		var sqlParentID sql.NullInt64
		var targetNamespace string
		var autoDeploy bool

		err = rows.Scan(&stepID, &stepNumber, &sqlParentID, &targetNamespace, &autoDeploy)
		if err != nil {
			return nil, err
		}
		if sqlParentID.Valid {
			// in db parent_id is an int, so should be safe
			parentID = int(sqlParentID.Int64)
		}
		step := &engine.PipelineStep{
			ID:              stepID,
			StepNumber:      stepNumber,
			ParentID:        parentID,
			PipelineID:      pipelineID,
			TargetNamespace: targetNamespace,
			AutomaticDeploy: autoDeploy,
			NextSteps:       []*engine.PipelineStep{},
		}
		stepsMap[step.ID] = step
		// Add itself to list of nextsteps of parent step, if any
		if step.ParentID > 0 {
			if parent, found := stepsMap[step.ParentID]; found {
				parent.NextSteps = append(parent.NextSteps, step)
			}
		}
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	//select steps which are root of the tree
	roots := []*engine.PipelineStep{}
	for _, s := range stepsMap {
		if s.ParentID == 0 {
			roots = append(roots, s)
		}
	}

	deployment := &engine.Deployment{
		ID:            id,
		ChartName:     chart,
		RepositoryURL: repositoryURL,
		CreationDate:  creationDate,
		LastUpdate:    lastUpdate,
		Pipeline: &engine.Pipeline{
			ID:         pipelineID,
			LastUpdate: pipelineLastUpdate,
			RootSteps:  roots,
		},
		Releases: releases,
	}
	return deployment, nil
}
