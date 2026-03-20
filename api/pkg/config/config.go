// carrega e valida todas as configs

package config

import (
	"log"
	"os"
	"strings"
)

// Config é a struct que representa TODAS as configurações da aplicação.
// Em POO seria uma classe com atributos privados e getters.
type Config struct {
	Port        string
	GinMode     string
	AllowedOrigins []string // quais domínios podem chamar a API
	TrustedProxies []string // whitelist Cloudflare
	GitHubWebhookSecret string
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


	return &Config{
		Port:    getEnvOrDefault("PORT", "8080"),
		GinMode: getEnvOrDefault("GIN_MODE", "release"),

		// Divide "https://app.com,https://dash.com" em uma lista
		AllowedOrigins: trimOrigins(strings.Split(origins, ",")),

		GitHubWebhookSecret: secret, 

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