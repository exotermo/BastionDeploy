package handler

import (
	"exodeploy/internal/domain"
	"exodeploy/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *gin.Engine, webhookSecret string, deployRepo domain.DeployRepository, rdb *redis.Client) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": "BastionDeploy API", "version": "1.0.0"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", HealthCheck)
		v1.GET("/stats", NewStatsHandler(deployRepo))
		v1.GET("/apps/status", NewAppsStatusHandler(deployRepo))
		v1.GET("/deploys", NewDeploysHandler(deployRepo))

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