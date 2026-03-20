package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebhookHandler recebe o payload do GitHub/GitLab
func WebhookHandler(c *gin.Context) {

	/*
	Tenta transformar o corpo da requisição (JSON) numa variável Go. 
	O &payload passa o endereço de memória da variável 
	isso é um ponteiro, algo que o Python esconde de você mas o Go expõe. 
	Basicamente diz: "preenche essa variável aqui".
	*/

	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}

	// TODO: chamar o service de deploy
	c.JSON(http.StatusAccepted, gin.H{
		"message": "webhook recebido, deploy enfileirado",
	})


}

// StatusHandler retorna o estado de uma app específica
func StatusHandler(c *gin.Context) {

	/*
	Pega o valor dinâmico da URL — se a rota for `/status/meu-bot`, o `appName` vai ser `"meu-bot"`.
	*/
	appName := c.Param("app")

	

	// TODO: buscar status real no banco
	c.JSON(http.StatusOK, gin.H{
		"app":    appName,
		"status": "running",
	})
}
