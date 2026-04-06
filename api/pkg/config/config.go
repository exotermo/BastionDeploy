// carrega e valida todas as configs

package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Config é a struct que representa TODAS as configurações da aplicação.
// Em POO seria uma classe com atributos privados e getters.
type Config struct {
	Port                string
	GinMode             string
	AllowedOrigins      []string // quais domínios podem chamar a API
	TrustedProxies      []string // whitelist Cloudflare
	GitHubWebhookSecret string
	APIKey              string // X-API-Key — protege os endpoints da API
	// Banco de dados
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	// Redis
	RedisAddr     string
	RedisPassword string
}

// Load lê as variáveis de ambiente e devolve uma Config preenchida.
// É o nosso "construtor" — equivale ao __init__ do Python.
func Load() *Config {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		log.Fatal("ALLOWED_ORIGINS não definido no .env — obrigatório por segurança")
	}

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("GITHUB_WEBHOOK_SECRET não definido — obrigatório por segurança")
	}

	secret = strings.TrimSpace(secret)

	cfg := &Config{
		Port:    getEnvOrDefault("PORT", "8080"),
		GinMode: getEnvOrDefault("GIN_MODE", "release"),

		// Divide "https://app.com,https://dash.com" em uma lista
		AllowedOrigins: trimOrigins(strings.Split(origins, ",")),

		GitHubWebhookSecret: secret,
		APIKey:              strings.TrimSpace(os.Getenv("API_KEY")),

		// Banco de dados -> postgres 15
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "exodeploy"),
		DBPassword: strings.TrimSpace(os.Getenv("DB_PASSWORD")),
		DBName:     getEnvOrDefault("DB_NAME", "exodeploy"),
		DBSSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),

		// IPs do Cloudflare (proxy reverso deles)
		// Lista completa: https://www.cloudflare.com/ips/
		TrustedProxies: []string{
			"173.245.48.0/20",
			"103.21.244.0/22",
			"103.22.200.0/22",
			"103.31.4.0/22",
			"141.101.64.0/18",
			"108.162.192.0/18",
			"190.93.240.0/20",
			"188.114.96.0/20",
			"197.234.240.0/22",
			"198.41.128.0/17",
			"162.158.0.0/15",
			"104.16.0.0/13",
			"104.24.0.0/14",
			"172.64.0.0/13",
			"131.0.72.0/22",
		},
		RedisAddr:     getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword: strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
	}

	cfg.Validate()
	return cfg
}

// Validate checa configurações críticas e loga avisos.
func (c *Config) Validate() {
	if c.APIKey == "" {
		log.Println("⚠️  API_KEY não definida — endpoints ficam sem autenticação (inseguro para produção)")
	}
}

// getEnvOrDefault é um helper privado (letra minúscula = privado em Go,
// igual métodos privados em Java/Kotlin)
func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// trimOrigins remove espaços em branco de cada origem
// evita erro se o .env tiver "http://localhost:3000, https://exotermo.dev" (com espaço)
func trimOrigins(origins []string) []string {
	result := make([]string, len(origins))
	for i, o := range origins {
		result[i] = strings.TrimSpace(o)
	}
	return result
}

// DatabaseURL retorna a connection string pro PostgreSQL
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}
