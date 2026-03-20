package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
*gin.Context é o objeto que contém tudo da requisição HTTP 
c aqui é igual ao req/res do Express, ou o request do FastAPI.
*/
func HealthCheck(c *gin.Context) {

	/*
	Responde com JSON. http.StatusOK = código 200. gin.H{} é um atalho pra map[string]any{} — basicamente um dicionário Go.
	*/
	c.JSON(http.StatusOK, gin.H{
		"status":    "online",
		"service":   "ExoDeploy API",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}