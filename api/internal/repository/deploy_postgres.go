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


func (r *PostgresDeployRepository) GetStats() (*domain.Stats, error) {
	stats := &domain.Stats{}

	// Total de deploys
	r.db.QueryRow(`SELECT COUNT(*) FROM deploys`).Scan(&stats.TotalDeploys)

	// Taxa de sucesso
	var success int
	r.db.QueryRow(`SELECT COUNT(*) FROM deploys WHERE status = 'success'`).Scan(&success)
	if stats.TotalDeploys > 0 {
		stats.SuccessRate = float64(success) / float64(stats.TotalDeploys) * 100
	}

	// Apps únicas
	r.db.QueryRow(`SELECT COUNT(DISTINCT app_name) FROM deploys`).Scan(&stats.ActiveApps)

	// Último deploy
	var lastDeploy string
	err := r.db.QueryRow(`
		SELECT created_at::text FROM deploys 
		ORDER BY created_at DESC LIMIT 1
	`).Scan(&lastDeploy)
	if err == nil {
		stats.LastDeployAt = &lastDeploy
	}

	return stats, nil
}

func (r *PostgresDeployRepository) GetAppsStatus() ([]*domain.AppStatus, error) {
	rows, err := r.db.Query(`
		SELECT 
			app_name,
			(SELECT status FROM deploys d2 
			 WHERE d2.app_name = d1.app_name 
			 ORDER BY created_at DESC LIMIT 1) as last_status
		FROM deploys d1
		GROUP BY app_name
		ORDER BY app_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*domain.AppStatus
	for rows.Next() {
		app := &domain.AppStatus{}
		var lastStatus string
		rows.Scan(&app.Name, &lastStatus)

		if lastStatus == "success" || lastStatus == "running" {
			app.Status = "UP"
		} else {
			app.Status = "DOWN"
		}
		app.Uptime = "—" // implementar depois com métricas reais
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *PostgresDeployRepository) GetRecent(limit int) ([]*domain.Deploy, error) {
	rows, err := r.db.Query(`
		SELECT id, app_name, branch, commit_sha, status, triggered_by, created_at, updated_at
		FROM deploys
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deploys []*domain.Deploy
	for rows.Next() {
		d := &domain.Deploy{}
		rows.Scan(&d.ID, &d.AppName, &d.Branch, &d.CommitSHA,
			&d.Status, &d.TriggeredBy, &d.CreatedAt, &d.UpdatedAt)
		deploys = append(deploys, d)
	}
	return deploys, nil
}