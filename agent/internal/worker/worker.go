package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"

	"exodeploy/agent/internal/notifier"
	"exodeploy/agent/internal/provisioner"
	"exodeploy/agent/internal/runner"
)

const (
	QueueName = "bastiondeploy:jobs" // chave no Redis
	PollInterval = 3 * time.Second   // verifica a fila a cada 3s
)

// Job representa o payload que a API publica no Redis
type Job struct {
	DeployID      string `json:"deploy_id"`
	AppName       string `json:"app_name"`
	Branch        string `json:"branch"`
	CommitSHA     string `json:"commit_sha"`
	RepoURL       string `json:"repo_url"`
	TriggeredBy   string `json:"triggered_by"`
	Domain        string `json:"domain"`
	EnableSSL     bool   `json:"enable_ssl"`
	ContainerPort int    `json:"container_port"`
}

type Worker struct {
	redis         *redis.Client
	db            *sql.DB
	runner        *runner.DockerRunner
	notifier      *notifier.Discord
	provisioner   *provisioner.NginxProvisioner
	tunnelMode    bool
}

func New(cfg *Config) (*Worker, error) {
	// Conecta no Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis inacessível: %w", err)
	}
	log.Println("✅ Redis conectado")

	// Conecta no PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("PostgreSQL inacessível: %w", err)
	}
	log.Println("✅ PostgreSQL conectado")

	w := &Worker{
		redis:       rdb,
		db:          db,
		runner:      runner.NewDockerRunner(),
		notifier:    notifier.NewDiscord(cfg.DiscordWebhookURL),
		provisioner: provisioner.NewNginxProvisioner(cfg.CertbotEmail),
		tunnelMode:  cfg.UseTunnelMode,
	}

	// Setup Cloudflare Tunnel se habilitado
	if cfg.UseTunnelMode {
		cf := provisioner.NewCloudflareProvisioner(cfg.CloudflareTunnelToken)
		if err := cf.SetupTunnel(); err != nil {
			log.Printf("⚠️  Cloudflare tunnel falhou: %v (deploy continuará sem tunnel)", err)
		}
	}

	return w, nil
}

func (w *Worker) Close() {
	w.redis.Close()
	w.db.Close()
}

// Start fica em loop consumindo jobs da fila até o contexto cancelar
func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.processNext(ctx)
		}
	}
}

func (w *Worker) processNext(ctx context.Context) {
	// BRPOP — bloqueia até chegar um job ou timeout de 3s
	// Assim o agent não fica consumindo CPU em loop vazio
	result, err := w.redis.BRPop(ctx, PollInterval, QueueName).Result()
	if err == redis.Nil {
		return // timeout, fila vazia — tenta de novo
	}
	if err != nil {
		log.Printf("❌ Erro ao ler fila: %v", err)
		return
	}

	// result[0] = nome da fila, result[1] = payload
	var job Job
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		log.Printf("❌ Payload inválido: %v", err)
		return
	}

	log.Printf("🚀 Processando deploy [%s] app=%s branch=%s", job.DeployID, job.AppName, job.Branch)

	// Atualiza status para "running"
	w.updateStatus(job.DeployID, "running")

	// Executa o deploy (container)
	err = w.runner.Deploy(ctx, job.AppName, job.RepoURL, job.Branch, &runner.DeployOptions{
		Domain:        job.Domain,
		EnableSSL:     job.EnableSSL,
		ContainerPort: job.ContainerPort,
	})
	if err != nil {
		log.Printf("❌ Deploy falhou [%s]: %v", job.DeployID, err)
		w.updateStatus(job.DeployID, "failed")
		w.notifier.Send(fmt.Sprintf(
			"❌ **Deploy falhou**\nApp: `%s`\nBranch: `%s`\nCommit: `%s`\nErro: `%v`",
			job.AppName, job.Branch, job.CommitSHA, err,
		))
		return
	}

	// Provisiona Nginx/Certbot/Tunnel se domínio foi informado
	if job.Domain != "" {
		tunnelMode := w.tunnelMode
		if err := w.provisioner.Provision(provisioner.NginxConfig{
			AppName:       job.AppName,
			Domain:        job.Domain,
			ContainerPort: job.ContainerPort,
			EnableSSL:     job.EnableSSL,
			TunnelMode:    tunnelMode,
		}); err != nil {
			log.Printf("⚠️  Nginx falhou para %s: %v", job.Domain, err)
			// Não falha o deploy — container está rodando, SSL/HTTP pode ser retried
		}
	}

	log.Printf("✅ Deploy concluído [%s]", job.DeployID)
	w.updateStatus(job.DeployID, "success")
	w.notifier.Send(fmt.Sprintf(
		"✅ **Deploy concluído**\nApp: `%s`\nBranch: `%s`\nCommit: `%s`\nPor: `%s`",
		job.AppName, job.Branch, job.CommitSHA, job.TriggeredBy,
	))
}

func (w *Worker) updateStatus(deployID, status string) {
	_, err := w.db.Exec(
		`UPDATE deploys SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, deployID,
	)
	if err != nil {
		log.Printf("❌ Erro ao atualizar status: %v", err)
	}
}