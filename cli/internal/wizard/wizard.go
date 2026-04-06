// Package wizard oferece um assistente interativo de configuração
package wizard

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

	"exodeploy/cli/internal/config"
)

type Wizard struct {
	cfg   *config.EnvConfig
	scanner *bufio.Scanner
}

func New(cfg *config.EnvConfig) *Wizard {
	return &Wizard{
		cfg:     cfg,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (w *Wizard) RunFull() error {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║        ExoDeploy — Setup Wizard              ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()

	if err := w.configureDatabase(); err != nil {
		return err
	}
	if err := w.configureDiscord(); err != nil {
		return err
	}
	if err := w.configureWebhook(); err != nil {
		return err
	}
	if err := w.configureExposeMode(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✅ Configurações salvas!")
	fmt.Println()
	fmt.Println("Próximos passos:")
	fmt.Println("  1. cd api && go run cmd/api/main.go    (ou docker compose up -d)")
	fmt.Println("  2. cd agent && go run cmd/agent/main.go")
	fmt.Println("  3. Teste com: curl http://localhost:8080")
	return nil
}

// --- steps ---

func (w *Wizard) configureDatabase() error {
	fmt.Println("────── Database & Redis ──────")
	defaultHost    := w.readOr("DB host", "127.0.0.1")
	defaultPort    := w.readOr("DB port", "5432")
	defaultUser    := w.readOr("DB user", "exodeploy")
	defaultPass    := w.readSecret("DB password")
	defaultName    := w.readOr("DB name", "exodeploy")

	if defaultPass == "" {
		return fmt.Errorf("senha do banco é obrigatória")
	}

	if err := w.setAll("api", map[string]string{
		"DB_HOST":     defaultHost,
		"DB_PORT":     defaultPort,
		"DB_USER":     defaultUser,
		"DB_PASSWORD": defaultPass,
		"DB_NAME":     defaultName,
		"DB_SSLMODE":  "disable",
	}); err != nil {
		return err
	}
	if err := w.setAll("agent", map[string]string{
		"DB_HOST":     defaultHost,
		"DB_PORT":     defaultPort,
		"DB_USER":     defaultUser,
		"DB_PASSWORD": defaultPass,
		"DB_NAME":     defaultName,
		"DB_SSLMODE":  "disable",
	}); err != nil {
		return err
	}

	fmt.Println()

	// Redis
	fmt.Println("────── Redis ──────")
	redisAddr := w.readOr("Redis address", "127.0.0.1:6379")
	redisPass := w.readSecret("Redis password (vazio = sem senha)")

	if err := w.setAll("api", map[string]string{
		"REDIS_ADDR":     redisAddr,
		"REDIS_PASSWORD": redisPass,
	}); err != nil {
		return err
	}
	if err := w.setAll("agent", map[string]string{
		"REDIS_ADDR":     redisAddr,
		"REDIS_PASSWORD": redisPass,
	}); err != nil {
		return err
	}

	fmt.Println()
	return nil
}

func (w *Wizard) configureDiscord() error {
	fmt.Println("────── Discord Notifications ──────")
	fmt.Println("Cole a URL do Webhook do Discord (ou vazio para pular):")

	url := w.readLine("Discord webhook URL → ")

	// Validação rápida
	if url != "" && !strings.Contains(url, "discord.com") && !strings.Contains(url, "discordapp.com") {
		fmt.Println("⚠️  URL parece inválida — não contém 'discord.com'. Continuar assim mesmo? [n]")
		if w.readConfirmDefaultNo() {
			return nil
		}
	}

	return w.setAll("agent", map[string]string{
		"DISCORD_WEBHOOK_URL": url,
	})
}

func (w *Wizard) configureWebhook() error {
	fmt.Println("\n────── GitHub Webhook ──────")
	fmt.Println("Gere um secret aleatório e configure no GitHub > Settings > Webhooks")
	secret := w.readLine("GitHub webhook secret → ")
	if secret == "" {
		secret = GenerateSecret()
		fmt.Printf("Nenhum valor informado, gerando: %s\n", secret)
	}

	return w.setAll("api", map[string]string{
		"GITHUB_WEBHOOK_SECRET": secret,
		"ALLOWED_ORIGINS":       "http://localhost:5173,http://localhost:3000",
		"PORT":                  "8080",
		"GIN_MODE":              "release",
	})
}

func (w *Wizard) configureExposeMode() error {
	fmt.Println("\n────── Expose Mode ──────")
	fmt.Println("Como seus apps serão acessados?")
	fmt.Println("  1) Cloudflare Tunnel (recomendado — zero portas abertas na VPS)")
	fmt.Println("  2) Nginx direto na máquina (portas 80/443 abertas)")
	fmt.Println()

	option := w.readLine("> Escolha [1/2] → ")

	if option == "1" || option == "" {
		fmt.Println("\n☁️  Modo: Cloudflare Tunnel")
		fmt.Println("\nPara obter o token:")
		fmt.Println("  1. Acesse https://one.dash.cloudflare.com/")
		fmt.Println("  2. Networks > Tunnels > Create a tunnel")
		fmt.Println("  3. Copie o token da seção 'Install and run a connector'")
		fmt.Println()

		token := w.readLine("Cloudflare Tunnel Token → ")

		if err := w.setAll("agent", map[string]string{
			"USE_TUNNEL_MODE":         "true",
			"CLOUDFLARE_TUNNEL_TOKEN": token,
			"CERTBOT_EMAIL":           "",
		}); err != nil {
			return err
		}
	} else {
		fmt.Println("\n🖥️  Modo: Nginx local (portas 80/443 abertas)")
		fmt.Println("O Certbot será usado para gerar certificados Let's Encrypt.")
		fmt.Println()

		email := w.readLine("Email para Certbot (Let's Encrypt) → ")

		if err := w.setAll("agent", map[string]string{
			"USE_TUNNEL_MODE":         "false",
			"CLOUDFLARE_TUNNEL_TOKEN": "",
			"CERTBOT_EMAIL":           email,
		}); err != nil {
			return err
		}
	}
	return nil
}

// --- helpers ---

func (w *Wizard) readLine(prompt string) string {
	fmt.Print(prompt)
	if !w.scanner.Scan() {
		return ""
	}
	return strings.TrimSpace(w.scanner.Text())
}

// readOr mostra o valor atual como default; se o usuário digitar vazio, retorna o default
func (w *Wizard) readOr(label, defaultVal string) string {
	current, err := w.cfg.Get("agent", label)
	if err != nil {
		current = ""
	}
	display := current
	if display == "" {
		display = defaultVal
	}

	val := w.readLine(fmt.Sprintf("%s [%s] → ", label, display))
	if val == "" {
		val = current
		if val == "" {
			val = defaultVal
		}
	}
	return val
}

func (w *Wizard) readSecret(prompt string) string {
	fmt.Print(prompt + " → ")
	// read normal text sem prompt de eco (fallback: lê normal)
	// Go padrão não suporta leitura sem echo sem cgo; por simplicidade lê normal
	if !w.scanner.Scan() {
		return ""
	}
	return strings.TrimSpace(w.scanner.Text())
}

func (w *Wizard) readConfirmDefaultNo() bool {
	fmt.Print("Continuar mesmo assim? [y/N] → ")
	if w.scanner.Scan() {
		ans := strings.TrimSpace(w.scanner.Text())
		return ans == "y" || ans == "Y" || ans == "yes"
	}
	return false
}

func (w *Wizard) setAll(component string, kv map[string]string) error {
	for k, v := range kv {
		if err := w.cfg.Set(component, k, v); err != nil {
			return fmt.Errorf("erro ao gravar %s em %s: %w", k, component, err)
		}
	}
	return nil
}

func GenerateSecret() string {
	b := make([]byte, 32)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return "ghsec_" + string(b)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
