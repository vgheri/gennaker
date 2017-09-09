package pg

import (
	"database/sql"

	"github.com/vgheri/gennaker/engine"
)

func createPipelineStep(tx *sql.Tx, deploymentID int, step *engine.PipelineStep) error {
	if step == nil {
		return engine.ErrInvalidPipeline
	}
	query := `INSERT INTO pipeline_step(step_number, parent_step_number, deployment_id, target_namespace, auto_deploy)
  VALUES($1, $2, $3, $4, $5) RETURNING id;`
	var id int
	var row *sql.Row
	if step.ParentStepNumber == 0 {
		row = tx.QueryRow(query, step.StepNumber, nil, deploymentID,
			step.TargetNamespace, step.AutomaticDeploy)

	} else {
		row = tx.QueryRow(query, step.StepNumber, step.ParentStepNumber, deploymentID,
			step.TargetNamespace, step.AutomaticDeploy)
	}
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	step.ID = id
	step.DeploymentID = deploymentID
	for _, cs := range step.NextSteps {
		err = createPipelineStep(tx, deploymentID, cs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *pgRepository) getDeploymentPipeline(deploymentID int) ([]*engine.PipelineStep, error) {
	// Build the pipeline
	query := `SELECT id, step_number, parent_step_number, target_namespace, auto_deploy
  FROM pipeline_step
  WHERE deployment_id = $1
  ORDER BY step_number asc;`
	rows, err := r.db.Query(query, deploymentID)
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
			DeploymentID:     deploymentID,
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
	return roots, nil
}
