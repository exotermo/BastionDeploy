package handler

import (
	"exodeploy/internal/domain"
	"exodeploy/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *gin.Engine, webhookSecret string, deployRepo domain.DeployRepository, rdb *redis.Client) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", HealthCheck)

		deploy := v1.Group("/deploy")
		{
			deploy.POST("/webhook",
				middleware.ValidateGitHubWebhook(webhookSecret),
				NewWebhookHandler(deployRepo, rdb),
			)
			deploy.GET("/status/:app", NewStatusHandler(deployRepo))
		}
	}
}