package domain

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

type Deploy struct {
	ID          string    `json:"id"`
	AppName     string    `json:"app_name"`
	Branch      string    `json:"branch"`
	CommitSHA   string    `json:"commit_sha"`
	Status      Status    `json:"status"`
	TriggeredBy string    `json:"triggered_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Stats — dados para os cards do dashboard
type Stats struct {
	TotalDeploys  int     `json:"total_deploys"`
	SuccessRate   float64 `json:"success_rate"`
	ActiveApps    int     `json:"active_apps"`
	LastDeployAt  *string `json:"last_deploy_at"`
}

// AppStatus — status de cada app para o painel lateral
type AppStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "UP" ou "DOWN"
	Uptime string `json:"uptime"`
}

type DeployRepository interface {
	Save(deploy *Deploy) error
	FindByID(id string) (*Deploy, error)
	FindByApp(appName string) ([]*Deploy, error)
	UpdateStatus(id string, status Status) error
	// Novos métodos
	GetStats() (*Stats, error)
	GetAppsStatus() ([]*AppStatus, error)
	GetRecent(limit int) ([]*Deploy, error)
}