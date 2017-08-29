package pg

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vgheri/gennaker/engine"
)

func (r *pgRepository) CreateRelease(release *engine.Release) (int, error) {
	var releaseID int
	var values, chartVersion sql.NullString
	if len(strings.TrimSpace(release.Values)) != 0 {
		values.Valid = true
		values.String = release.Values
	}
	if len(strings.TrimSpace(release.ChartVersion)) != 0 {
		chartVersion.Valid = true
		chartVersion.String = release.ChartVersion
	}

	query := `INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, chart_version, status)
  VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	tx, err := r.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "Cannot init transaction")
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, release.Name, release.DeploymentID, release.ImageTag, release.Namespace,
		values, release.Chart, chartVersion, release.Status).Scan(&releaseID)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return 0, errors.Wrap(err, "Cannot insert release")
	}
	if err = tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "Cannot commit transaction")
	}
	return releaseID, nil
}
