// Package config lê e grava os arquivos .env da API e do Agent
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EnvConfig struct {
	InstallDir string // onde está a instalação do ExoDeploy
}

func New(installDir string) *EnvConfig {
	return &EnvConfig{InstallDir: installDir}
}

// ResolveInstallDir detecta a instalação do ExoDeploy
// Busca em ordem: --dir flag, EXODEPLOY_DIR env, cwd se contiver api/agent/
func ResolveInstallDir(explicit string) (string, error) {
	if explicit != "" {
		if err := validateInstallDir(explicit); err == nil {
			return explicit, nil
		}
		return "", fmt.Errorf("diretório inválido: %s", explicit)
	}
	if dir := os.Getenv("EXODEPLOY_DIR"); dir != "" {
		return dir, nil
	}
	cwd, _ := os.Getwd()
	if err := validateInstallDir(cwd); err == nil {
		return cwd, nil
	}
	return "", fmt.Errorf("não foi possível encontrar o ExoDeploy — use --dir <path>")
}

func validateInstallDir(dir string) error {
	if _, err := os.Stat(filepath.Join(dir, "api")); err != nil {
		return fmt.Errorf("pasta 'api' não encontrada")
	}
	if _, err := os.Stat(filepath.Join(dir, "agent")); err != nil {
		return fmt.Errorf("pasta 'agent' não encontrada")
	}
	return nil
}

// readEnv lê um arquivo .env e retorna como map
func ReadEnvFile(path string) (map[string]string, error) {
	result := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok {
			result[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return result, scanner.Err()
}

// WriteEnvFile grava um map como .env
func WriteEnvFile(path string, values map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "# ExoDeploy .env — gerado por exoctl")
	fmt.Fprintln(w)

	// Grava na ordem para facilitar leitura humana
	order := []string{
		"PORT", "GIN_MODE", "ALLOWED_ORIGINS", "GITHUB_WEBHOOK_SECRET",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"REDIS_ADDR", "REDIS_PASSWORD",
		"DISCORD_WEBHOOK_URL",
		"CERTBOT_EMAIL",
		"CLOUDFLARE_TUNNEL_TOKEN", "USE_TUNNEL_MODE",
	}
	written := make(map[string]bool)
	for _, key := range order {
		if v, ok := values[key]; ok {
			fmt.Fprintf(w, "%s=%s\n", key, v)
			written[key] = true
		}
	}
	// Chaves não-ordered restantes
	for k, v := range values {
		if !written[k] {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}
	return w.Flush()
}

// Get lê um valor do .env da API
func (e *EnvConfig) Get(component, key string) (string, error) {
	path := e.envPath(component)
	values, err := ReadEnvFile(path)
	if err != nil {
		return "", err
	}
	return values[key], nil
}

func (e *EnvConfig) Set(component, key, value string) error {
	path := e.envPath(component)
	values := make(map[string]string)

	// Tenta ler, senão começa vazio
	if f, err := os.Open(path); err == nil {
		f.Close()
		values, _ = ReadEnvFile(path)
	}

	values[key] = value
	return WriteEnvFile(path, values)
}

func (e *EnvConfig) envPath(component string) string {
	switch component {
	case "api":
		return filepath.Join(e.InstallDir, "api", ".env")
	case "agent":
		return filepath.Join(e.InstallDir, "agent", ".env")
	case "dashboard":
		return filepath.Join(e.InstallDir, "dashboard", ".env")
	default:
		return filepath.Join(e.InstallDir, ".env")
	}
}

func MustSet(values map[string]string, key string, val string) {
	values[key] = val
}

func (e *EnvConfig) Has(component, key string) (bool, error) {
	v, err := e.Get(component, key)
	return v != "", err
}

func StrBool(s string) bool {
	sl := strings.ToLower(strings.TrimSpace(s))
	return sl == "true" || sl == "yes" || sl == "1"
}

func BoolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func StrInt(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
