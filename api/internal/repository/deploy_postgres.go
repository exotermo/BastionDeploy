package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"exodeploy/internal/domain"
)

type PostgresDeployRepository struct {
	db *sql.DB
}

func NewPostgresDeployRepository(db *sql.DB) domain.DeployRepository {
	return &PostgresDeployRepository{db: db}
}

func (r *PostgresDeployRepository) Save(deploy *domain.Deploy) error {
	query := `
		INSERT INTO deploys (app_name, branch, commit_sha, status, triggered_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query,
		deploy.AppName,
		deploy.Branch,
		deploy.CommitSHA,
		deploy.Status,
		deploy.TriggeredBy,
	).Scan(&deploy.ID, &deploy.CreatedAt, &deploy.UpdatedAt)
}

func (r *PostgresDeployRepository) FindByID(id string) (*domain.Deploy, error) {
	query := `
		SELECT id, app_name, branch, commit_sha, status, triggered_by, created_at, updated_at
		FROM deploys WHERE id = $1
	`
	deploy := &domain.Deploy{}
	err := r.db.QueryRow(query, id).Scan(
		&deploy.ID, &deploy.AppName, &deploy.Branch, &deploy.CommitSHA,
		&deploy.Status, &deploy.TriggeredBy, &deploy.CreatedAt, &deploy.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("deploy %s não encontrado", id)
	}
	return deploy, err
}

func (r *PostgresDeployRepository) FindByApp(appName string) ([]*domain.Deploy, error) {
	query := `
		SELECT id, app_name, branch, commit_sha, status, triggered_by, created_at, updated_at
		FROM deploys WHERE app_name = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, appName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deploys []*domain.Deploy
	for rows.Next() {
		d := &domain.Deploy{}
		err := rows.Scan(
			&d.ID, &d.AppName, &d.Branch, &d.CommitSHA,
			&d.Status, &d.TriggeredBy, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		deploys = append(deploys, d)
	}
	return deploys, nil
}

func (r *PostgresDeployRepository) UpdateStatus(id string, status domain.Status) error {
	query := `UPDATE deploys SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}