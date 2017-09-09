package pg

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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

	query := `INSERT INTO release(name, deployment_id, image_tag, namespace, values, chart, chart_version, revision, status)
  VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	tx, err := r.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "Cannot init transaction")
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, release.Name, release.DeploymentID, release.ImageTag, release.Namespace,
		values, release.Chart, chartVersion, release.Revision, release.Status).Scan(&releaseID)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return 0, errors.Wrap(err, "Cannot insert release")
	}
	if err = tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "Cannot commit transaction")
	}
	return releaseID, nil
}

func (r *pgRepository) GetDeploymentReleases(deploymentID int) ([]*engine.Release, error) {
	query := `SELECT id, name, image_tag, timestamp, namespace, values, chart,
	chart_version, revision, status
	FROM release
	WHERE deployment_id = $1
	ORDER BY timestamp desc;`
	rows, err := r.db.Query(query, deploymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	releases := []*engine.Release{}
	for rows.Next() { // TODO: replace with call to getRelease
		var releaseID, revision int
		var timestamp time.Time
		var imageTag, namespace, chart, name string
		var values, chartVersion sql.NullString
		var status uint8
		err = rows.Scan(&releaseID, &name, &imageTag, &timestamp, &namespace,
			&values, &chart, &chartVersion, &revision, &status)
		if err != nil {
			return nil, err
		}
		release := &engine.Release{
			ID:           releaseID,
			Name:         name,
			ImageTag:     imageTag,
			DeploymentID: deploymentID,
			Date:         timestamp,
			Namespace:    namespace,
			Values:       values.String,
			Chart:        chart,
			ChartVersion: chartVersion.String,
			Revision:     revision,
			Status:       engine.GennakerReleaseOutcome(status),
		}
		releases = append(releases, release)
	}
	return releases, nil
}
