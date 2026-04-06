package worker

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	// Redis
	RedisAddr     string
	RedisPassword string

	// PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Discord
	DiscordWebhookURL string

	// Nginx/Certbot
	CertbotEmail string

	// Cloudflare Tunnel
	CloudflareTunnelToken string
	UseTunnelMode        bool

	// InstallDir — caminho do Agent para gerar scripts
	InstallDir string
}

func LoadConfig() *Config {
	return &Config{
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "exodeploy"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "exodeploy"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DiscordWebhookURL: getEnv("DISCORD_WEBHOOK_URL", ""),
		CertbotEmail:           getEnv("CERTBOT_EMAIL", ""),
		CloudflareTunnelToken:  getEnv("CLOUDFLARE_TUNNEL_TOKEN", ""),
		UseTunnelMode:          getEnv("USE_TUNNEL_MODE", "false") == "true",
		InstallDir:             detectInstallDir(),
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// detectInstallDir tenta encontrar a raiz do BastionDeploy
func detectInstallDir() string {
	exe, err := os.Executable()
	if err == nil {
		dir := exe
		// ../agent/cmd -> ../agent -> .. (raiz)
		for i := 0; i < 3; i++ {
			dir = dir[:len(dir)-len(filepath.Base(dir))-1]
			if _, err := os.Stat(filepath.Join(dir, "agent")); err == nil {
				return dir
			}
		}
	}
	cwd, _ := os.Getwd()
	return cwd
}