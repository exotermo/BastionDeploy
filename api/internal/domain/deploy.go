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

type DeployRepository interface {
	Save(deploy *Deploy) error
	FindByID(id string) (*Deploy, error)
	FindByApp(appName string) ([]*Deploy, error)
	UpdateStatus(id string, status Status) error
}