package worker

import (
	"fmt"
	"os"
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
		CertbotEmail:         getEnv("CERTBOT_EMAIL", ""),
		CloudflareTunnelToken: getEnv("CLOUDFLARE_TUNNEL_TOKEN", ""),
		UseTunnelMode:        getEnv("USE_TUNNEL_MODE", "false") == "true",
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