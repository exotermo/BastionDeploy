package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"exodeploy/internal/domain"
)

const QueueName = "bastiondeploy:jobs"

// Job é o payload publicado no Redis para o agent consumir
type Job struct {
	DeployID      string `json:"deploy_id"`
	AppName       string `json:"app_name"`
	Branch        string `json:"branch"`
	CommitSHA     string `json:"commit_sha"`
	RepoURL       string `json:"repo_url"`
	TriggeredBy   string `json:"triggered_by"`
	Domain        string `json:"domain"`
	EnableSSL     bool   `json:"enable_ssl"`
	ContainerPort int    `json:"container_port"`
}

func NewWebhookHandler(repo domain.DeployRepository, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Ref        string `json:"ref"`
			HeadCommit struct {
				ID string `json:"id"`
			} `json:"head_commit"`
			Repository struct {
				Name    string `json:"name"`
				CloneURL string `json:"clone_url"`
			} `json:"repository"`
			Pusher struct {
				Name string `json:"name"`
			} `json:"pusher"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
			return
		}

		// 1. Salva no PostgreSQL
		deploy := &domain.Deploy{
			AppName:     payload.Repository.Name,
			Branch:      payload.Ref,
			CommitSHA:   payload.HeadCommit.ID,
			Status:      domain.StatusPending,
			TriggeredBy: payload.Pusher.Name,
		}
		if err := repo.Save(deploy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar deploy"})
			return
		}

		// 2. Publica no Redis para o agent processar
		job := Job{
			DeployID:      deploy.ID,
			AppName:       deploy.AppName,
			Branch:        deploy.Branch,
			CommitSHA:     deploy.CommitSHA,
			RepoURL:       payload.Repository.CloneURL,
			TriggeredBy:   deploy.TriggeredBy,
			Domain:        c.GetHeader("X-Deploy-Domain"),       // ex: app.exemplo.com
			EnableSSL:     c.GetHeader("X-Deploy-SSL") == "true", // header opcional
			ContainerPort: 8080,                                  // default, sobrescrito pelo header
		}

		// Sobrescreve porta via header se fornecido
		if port := c.GetHeader("X-Deploy-Port"); port != "" {
			var p int
			if _, err := fmt.Sscanf(port, "%d", &p); err == nil && p > 0 {
				job.ContainerPort = p
			}
		}
		jobBytes, _ := json.Marshal(job)
		if err := rdb.LPush(context.Background(), QueueName, jobBytes).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao enfileirar deploy"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message":   "deploy enfileirado",
			"deploy_id": deploy.ID,
			"app":       deploy.AppName,
			"branch":    deploy.Branch,
		})
	}
}

func NewStatusHandler(repo domain.DeployRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		appName := c.Param("app")
		deploys, err := repo.FindByApp(appName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar deploys"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"app":     appName,
			"deploys": deploys,
		})
	}
}