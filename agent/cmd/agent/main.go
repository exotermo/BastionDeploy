package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"exodeploy/agent/internal/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum .env encontrado, usando variáveis do sistema")
	}

	cfg := worker.LoadConfig()

	log.Println("🤖 BastionDeploy Agent iniciado")
	log.Printf("📡 Conectando ao Redis em %s", cfg.RedisAddr)
	log.Printf("🗄  Conectando ao PostgreSQL em %s", cfg.DBHost)

	w, err := worker.New(cfg)
	if err != nil {
		log.Fatal("Erro ao iniciar worker:", err)
	}
	defer w.Close()

	// Contexto que cancela quando recebe SIGTERM ou SIGINT (Ctrl+C)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("⏳ Aguardando jobs na fila...")
	w.Start(ctx)
	log.Println("👋 Agent encerrado com sucesso")
}