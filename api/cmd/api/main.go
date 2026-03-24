package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"exodeploy/internal/handler"
	"exodeploy/internal/repository"
	"exodeploy/pkg/config"
	"exodeploy/pkg/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum .env encontrado, usando variáveis do sistema")
	}

	cfg := config.Load()

	// PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		log.Fatal("Erro ao abrir conexão com banco:", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("Banco inacessível:", err)
	}
	log.Println("✅ PostgreSQL conectado")

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Redis inacessível:", err)
	}
	log.Println("✅ Redis conectado")

	deployRepo := repository.NewPostgresDeployRepository(db)

	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.NewCORS(cfg.AllowedOrigins))

	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		log.Fatal("Erro ao configurar proxies:", err)
	}

	handler.RegisterRoutes(r, cfg.GitHubWebhookSecret, deployRepo, rdb)

	log.Printf("🚀 ExoDeploy API rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}