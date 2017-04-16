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
