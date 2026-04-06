// exoctl — CLI de gerenciamento do ExoDeploy
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"exodeploy/cli/internal/config"
	"exodeploy/cli/internal/wizard"
)

const version = "0.1.0"

var usage = `exoctl — ExoDeploy CLI Tool v` + version + `

Usage:
  exoctl <command> [options]

Commands:
  setup          Assistente interativo de configuração
  set            Define uma configuração específica
  get            Mostra uma configuração atual
  status         Resumo de todas as configurações
  gen-secret     Gera um GitHub webhook secret aleatório
  help           Mostra esta mensagem

Exemplos:
  exoctl setup
  exoctl set agent DISCORD_WEBHOOK_URL https://discord.com/api/webhooks/xxx
  exoctl set api ALLOWED_ORIGINS https://meuapp.com
  exoctl get agent USE_TUNNEL_MODE
  exoctl status
  exoctl set agent USE_TUNNEL_MODE yes    (atualiza modo de exposição)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	dirFlag := ""
	for _, a := range args {
		if strings.HasPrefix(a, "--dir=") {
			dirFlag = strings.TrimPrefix(a, "--dir=")
		}
	}

	switch cmd {
	case "help", "-h", "--help":
		fmt.Print(usage)
		return

	case "gen-secret":
		fmt.Println(wizard.GenerateSecret())
		return
	}

	installDir, err := config.ResolveInstallDir(dirFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	env := config.New(installDir)

	switch cmd {
	case "setup":
		w := wizard.New(env)
		if err := w.RunFull(); err != nil {
			fmt.Fprintf(os.Stderr, "\nErro: %v\n", err)
			os.Exit(1)
		}

	case "set":
		execSet(env, args)

	case "get":
		execGet(env, args)

	case "status":
		execStatus(env, installDir)

	default:
		fmt.Fprintf(os.Stderr, "Comando desconhecido: %s\n", cmd)
		fmt.Print(usage)
		os.Exit(1)
	}
}

func parseArgs(args []string) (string, string) {
	// filtra --dir=
	filtered := []string{}
	for _, a := range args {
		if !strings.HasPrefix(a, "--dir=") {
			filtered = append(filtered, a)
		}
	}
	// Espera: <component> <key> [value]
	if len(filtered) < 2 {
		fmt.Fprintf(os.Stderr, "Uso: exoctl set <component> <key> [value]\n")
		fmt.Fprintln(os.Stderr, "  component: api | agent | dashboard")
		os.Exit(1)
	}
	return filtered[0], filtered[1]
}

func execSet(env *config.EnvConfig, args []string) {
	filtered := []string{}
	for _, a := range args {
		if !strings.HasPrefix(a, "--dir=") {
			filtered = append(filtered, a)
		}
	}

	if len(filtered) < 3 {
		fmt.Fprintf(os.Stderr, "Uso: exoctl set <component> <key> <value>\n")
		os.Exit(1)
	}

	component, key := filtered[0], filtered[1]
	value := strings.Join(filtered[2:], " ")

	if err := env.Set(component, key, value); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	// Atalhos úteis com validação/auto-config
	switch {
	case key == "USE_TUNNEL_MODE":
		tunnel := config.StrBool(value)
		if tunnel {
			env.Set(component, "USE_TUNNEL_MODE", "true")
			fmt.Println("✅ Tunnel mode ativado — lembre-se de configurar CLOUDFLARE_TUNNEL_TOKEN")
		} else {
			env.Set(component, "USE_TUNNEL_MODE", "false")
			fmt.Println("✅ Tunnel mode desativado — Nginx vai expor portas 80/443 diretamente")
		}

	case key == "DISCORD_WEBHOOK_URL":
		if strings.Contains(value, "discord.com") || strings.Contains(value, "discordapp.com") {
			fmt.Println("✅ Discord webhook configurado")
		} else {
			fmt.Println("⚠️  URL pode não ser um webhook válido do Discord")
		}
	}

	fmt.Printf("✅ %s.%s = %s\n", component, key, value)
}

func execGet(env *config.EnvConfig, args []string) {
	component, key := parseArgs(args)
	val, err := env.Get(component, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler %s/%s: %v\n", component, key, err)
		os.Exit(1)
	}
	if val == "" {
		fmt.Println("(vazio)")
	} else {
		// Não mostrar senhas em texto puro
		if isSensitive(key) {
			fmt.Println(maskValue(val))
		} else {
			fmt.Println(val)
		}
	}
}

func execStatus(env *config.EnvConfig, installDir string) {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║        ExoDeploy — Status                    ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Printf("Install dir: %s\n\n", installDir)

	apiStatus(env)
	agentStatus(env)
}

func apiStatus(env *config.EnvConfig) {
	fmt.Println("── API ──")
	checkKey(env, "api", "PORT", false)
	checkKey(env, "api", "GIN_MODE", false)
	checkKey(env, "api", "ALLOWED_ORIGINS", false)
	checkKey(env, "api", "GITHUB_WEBHOOK_SECRET", true)
	checkKey(env, "api", "DB_HOST", false)
	checkKey(env, "api", "DB_PORT", false)
	checkKey(env, "api", "REDIS_ADDR", false)
	fmt.Println()
}

func agentStatus(env *config.EnvConfig) {
	fmt.Println("── Agent ──")
	checkKey(env, "agent", "DB_HOST", false)
	checkKey(env, "agent", "REDIS_ADDR", false)
	checkKey(env, "agent", "DISCORD_WEBHOOK_URL", true)
	checkKey(env, "agent", "USE_TUNNEL_MODE", false)
	checkKey(env, "agent", "CLOUDFLARE_TUNNEL_TOKEN", true)
	checkKey(env, "agent", "CERTBOT_EMAIL", true)
	fmt.Println()
}

func checkKey(env *config.EnvConfig, component, key string, sensitive bool) {
	val, err := env.Get(component, key)
	status := "❌ não configurado"
	if err == nil && val != "" {
		status = "✅ "
		if sensitive {
			status += maskValue(val)
		} else {
			status += val
		}
	}
	fmt.Printf("  %-30s %s\n", key, status)
}

func isSensitive(key string) bool {
	switch key {
	case "GITHUB_WEBHOOK_SECRET", "DISCORD_WEBHOOK_URL", "DB_PASSWORD", "REDIS_PASSWORD",
		"CLOUDFLARE_TUNNEL_TOKEN", "CERTBOT_EMAIL":
		return true
	}
	return false
}

func maskValue(s string) string {
	if len(s) <= 8 {
		return "•••••"
	}
	return s[:4] + "••••••••" + s[len(s)-4:]
}

func init() {
	// flag.Parse evita erros se algo parecer flag
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
}
