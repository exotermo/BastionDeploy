package handler // define o pacote, como uma classe


import (
	"github.com/gin-gonic/gin"
	"exodeploy/pkg/middleware"
)

// RegisterRoutes agora recebe o secret para passar ao middleware HMAC
func RegisterRoutes(r *gin.Engine, webhookSecret string) {
	
	
	// Grupo base da API v1

	/*
	agrupa todas as rotas com o prefixo /api/v1 
	evitamos boiderplates desnecessarios assim
	*/
	v1 := r.Group("/api/v1")

	{
		/*
		define quais funções respondem pra cada URL

		# nota mental pra lembrar como faz para parametrizar comentarios igual no kotlin 


		GET /api/v1/health -> chama HealthCheck
		POST /api/v1/deploy/webhook -> chama WebhookHandler
		GET /api/v1/deploy/status/meu-bot -> o :app é um parâmetro dinâmico (igual :id no Express.js)
		*/

		v1.GET("/health", HealthCheck)

		// Rotas de deploy
		deploy := v1.Group("/deploy")
		{
			// O middleware HMAC é aplicado somente aqui  
			// isso não afeta outras rotas da api 
			deploy.POST("/webhook",
				middleware.ValidateGitHubWebhook(webhookSecret),
				WebhookHandler,
			)
			deploy.GET("/status/:app", StatusHandler)
		}
	}
}