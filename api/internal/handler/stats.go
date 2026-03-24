package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"exodeploy/internal/domain"
)

func NewStatsHandler(repo domain.DeployRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := repo.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar stats"})
			return
		}
		c.JSON(http.StatusOK, stats)
	}
}

func NewAppsStatusHandler(repo domain.DeployRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apps, err := repo.GetAppsStatus()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar apps"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"apps": apps})
	}
}

func NewDeploysHandler(repo domain.DeployRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		deploys, err := repo.GetRecent(20)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar deploys"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deploys": deploys})
	}
}