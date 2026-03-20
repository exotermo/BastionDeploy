package main // ponto de entrada para o compilador 

/*
Go precisa de um package main 
diz pro compilador do Go esse arquivo é o programa principal
*/

import (
	"log"

	"github.com/gin-gonic/gin" // gin é o framework HTTP, 
	"github.com/joho/godotenv" // godotenv lê o .env

	"exodeploy/internal/handler" //  handler proprio do framework bastion

	"exodeploy/pkg/config"
	"exodeploy/pkg/middleware"
)

func main() { // inicializador principal


	// variáveis do .env 
	if err := godotenv.Load(); err != nil {

		log.Println("Nenhum .env encontrado, usando variáveis do sistema")

	}


	// valida todas as configs de uma vez
	cfg := config.Load()

	// define o modo do Gin (debug ou release)
	gin.SetMode(cfg.GinMode)

	r := gin.New()



	r.Use(gin.Logger())           // logger das requisições
	r.Use(gin.Recovery())         // panics evita crashar em runtime
	r.Use(middleware.SecurityHeaders()) // headers de segurança
	r.Use(middleware.NewCORS(cfg.AllowedOrigins)) // CORS configurado

	// diz pro Gin confiar somente nos proxies do Cloudflare
	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		log.Fatal("Erro ao configurar proxies:", err)
	}

	// secret para as rotas
	handler.RegisterRoutes(r, cfg.GitHubWebhookSecret)

	/*
	Debug inicial
	*/
	log.Printf("ExoDeploy API rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}