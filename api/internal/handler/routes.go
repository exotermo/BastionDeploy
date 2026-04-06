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
		// Público: load balancers e health checks
		v1.GET("/health", HealthCheck)

		// Protegidos por API Key
		authed := v1.Group("", middleware.RequireAPIKey())
		{
			authed.GET("/stats", NewStatsHandler(deployRepo))
			authed.GET("/apps/status", NewAppsStatusHandler(deployRepo))
			authed.GET("/deploys", NewDeploysHandler(deployRepo))

			deploy := authed.Group("/deploy")
			{
				deploy.GET("/status/:app", NewStatusHandler(deployRepo))
			}

			// Webhook: protegido por HMAC (checagem GitHub), não por API Key
			v1.POST("/deploy/webhook",
				middleware.ValidateGitHubWebhook(webhookSecret),
				NewWebhookHandler(deployRepo, rdb),
			)
		}
	}
}