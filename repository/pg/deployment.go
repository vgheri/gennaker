package pg

import (
	"database/sql"
	"time"

	"github.com/vgheri/gennaker/engine"
)

func (r *pgRepository) ListDeployments(limit, offset int) ([]*engine.Deployment, error) {
	query := `SELECT id, name, chart, repository_url, creation_date, last_update
  FROM deployment
  ORDER BY chart LIMIT $1 OFFSET $2;`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	deployments := []*engine.Deployment{}
	for rows.Next() {
		var id int
		var name, chart, repositoryURL string
		var creationDate, lastUpdate time.Time
		err = rows.Scan(&id, &name, &chart, &repositoryURL, &creationDate, &lastUpdate)
		if err != nil {
			return nil, err
		}
		deployment := &engine.Deployment{
			ID:            id,
			Name:          name,
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
	query := `SELECT name, chart, repository_url, creation_date, last_update
  FROM deployment
  WHERE id = $1`

	row := r.db.QueryRow(query, id)
	var name, chart, repositoryURL string
	var creationDate, lastUpdate time.Time
	err := row.Scan(&name, &chart, &repositoryURL, &creationDate,
		&lastUpdate)
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

	// Build the pipeline
	query = `SELECT id, step_number, parent_step_number, target_namespace, auto_deploy
  FROM pipeline_step
  WHERE deployment_id = $1
  ORDER BY step_number asc;`
	rows, err = r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stepsMap := make(map[int]*engine.PipelineStep)
	for rows.Next() {
		var stepID, stepNumber, parentStepNumber int
		var sqlParentStepNumber sql.NullInt64
		var targetNamespace string
		var autoDeploy bool

		err = rows.Scan(&stepID, &stepNumber, &sqlParentStepNumber, &targetNamespace, &autoDeploy)
		if err != nil {
			return nil, err
		}
		if sqlParentStepNumber.Valid {
			// in db parent_step_number is an int, so should be safe
			parentStepNumber = int(sqlParentStepNumber.Int64)
		}
		step := &engine.PipelineStep{
			ID:               stepID,
			StepNumber:       stepNumber,
			ParentStepNumber: parentStepNumber,
			DeploymentID:     id,
			TargetNamespace:  targetNamespace,
			AutomaticDeploy:  autoDeploy,
			NextSteps:        []*engine.PipelineStep{},
		}
		stepsMap[step.StepNumber] = step
		// Add itself to list of nextsteps of parent step, if any
		if step.ParentStepNumber > 0 {
			if parent, found := stepsMap[step.ParentStepNumber]; found {
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
		if s.ParentStepNumber == 0 {
			roots = append(roots, s)
		}
	}

	deployment := &engine.Deployment{
		ID:            id,
		Name:          name,
		ChartName:     chart,
		RepositoryURL: repositoryURL,
		CreationDate:  creationDate,
		LastUpdate:    lastUpdate,
		Pipeline:      roots,
		Releases:      releases,
	}
	return deployment, nil
}

func (r *pgRepository) CreateDeployment(deployment *engine.Deployment) error {
	if deployment == nil {
		return engine.ErrInvalidDeployment
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Create the deployment
	query := `INSERT INTO deployment(name, chart, repository_url)
	VALUES($1, $2, $3) RETURNING id, creation_date, last_update`
	row := r.db.QueryRow(query, deployment.Name, deployment.ChartName, deployment.RepositoryURL)
	var id int
	var creationDate, lastUpdate time.Time
	err = row.Scan(&id, &creationDate, &lastUpdate)
	if err != nil {
		return err
	}
	if deployment.Pipeline == nil || len(deployment.Pipeline) == 0 {
		return engine.ErrInvalidPipeline
	}
	for _, step := range deployment.Pipeline {
		err = createPipelineStep(tx, id, step)
		if err != nil {
			return err
		}
	}
	deployment.ID = id
	deployment.CreationDate = creationDate
	deployment.LastUpdate = lastUpdate
	return tx.Commit()
}
