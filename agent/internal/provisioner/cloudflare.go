package provisioner

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// CloudflareProvisioner gerencia o cloudflared tunnel na VPS
type CloudflareProvisioner struct {
	token string
	dir   string // dir de instalação do Agent (onde grava o script)
}

func NewCloudflareProvisioner(token, installDir string) *CloudflareProvisioner {
	return &CloudflareProvisioner{
		token: token,
		dir:   installDir,
	}
}

// SetupTunnel verifica se o tunnel está configurado e ativo.
// Se não estiver, gera um script de setup e retorna instruções.
// Isso permite que o Agent inicie sem precisar de sudo.
func (p *CloudflareProvisioner) SetupTunnel() error {
	if _, err := exec.LookPath("cloudflared"); err != nil {
		// Binário não instalado — instruções genéricas
		log.Println("⚠️  cloudflared não está instalado")
		p.printSetupInstructions()
		return fmt.Errorf("cloudflared não encontrado no PATH")
	}

	if p.token == "" {
		// Sem token — tunnel mode ativo sem setup
		log.Println("☁️  Tunnel mode ativo, mas sem token — rode setup_tunnel.sh")
		p.printSetupInstructions()
		return nil
	}

	// Verifica se o serviço já está rodando
	if p.isTunnelActive() {
		log.Println("☁️  Cloudflare tunnel está ativo")
		return nil
	}

	// Tunnel não ativo — gera script e instruções
	log.Println("⚠️  Cloudflare tunnel não está ativo")
	log.Println("   Gere e rode o script de setup:")
	log.Printf("   %s/setup_tunnel.sh\n", p.dir)

	if err := p.generateSetupScript(); err != nil {
		return fmt.Errorf("erro ao gerar script: %w", err)
	}

	log.Println("   Execute: sudo bash", filepath.Join(p.dir, "setup_tunnel.sh"))
	return nil
}

func (p *CloudflareProvisioner) isTunnelActive() bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "cloudflared")
	if cmd.Run() != nil {
		return false
	}

	// Verifica se o token no arquivo de configuração bate
	envFile := "/etc/default/cloudflared"
	content, err := os.ReadFile(envFile)
	if err != nil {
		return false
	}
	if !bytes.Contains(content, []byte("TUNNEL_TOKEN")) {
		return false
	}
	// Token presente no systemd service?
	serviceFile := "/etc/systemd/system/cloudflared.service"
	svcContent, err := os.ReadFile(serviceFile)
	if err != nil {
		return false
	}
	if !bytes.Contains(svcContent, []byte("ExecStart")) {
		return false
	}
	return true
}

func (p *CloudflareProvisioner) generateSetupScript() error {
	setupPath := filepath.Join(p.dir, "setup_tunnel.sh")

	// Já existe e tem token? — não sobrescreve
	if content, err := os.ReadFile(setupPath); err == nil {
		if bytes.Contains(content, []byte(p.token[:8])) {
			return nil // script válido já existe
		}
	}

	scriptTemplate := template.Must(template.New("setup").Parse(`#!/bin/bash
# BastionDeploy — Cloudflare Tunnel Setup Script
# Gerado automaticamente pelo Agent
# Execute com: sudo bash setup_tunnel.sh

set -e

CLOUDFLARE_BIN=$(command -v cloudflared)
if [ -z "$CLOUDFLARE_BIN" ]; then
  echo "❌ cloudflared não instalado. Instale antes."
  echo "   Ubuntu: sudo apt-get install cloudflared"
  echo "   Fedora: sudo dnf install cloudflared"
  exit 1
fi

echo "☁️  Configurando Cloudflare Tunnel..."

# Garante diretório
mkdir -p /etc/default

# Token
echo 'TUNNEL_TOKEN={{.Token}}' > /etc/default/cloudflared
chmod 600 /etc/default/cloudflared
echo "✅ Token configurado"

# Systemd service
cat > /etc/systemd/system/cloudflared.service << 'SVCEOF'
[Unit]
Description=Cloudflare Tunnel
After=network-online.target

[Service]
Type=simple
EnvironmentFile=/etc/default/cloudflared
ExecStart=/usr/bin/cloudflared tunnel --no-autoupdate run --token ${TUNNEL_TOKEN}
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SVCEOF

echo "✅ Serviço systemd criado"

# Detecta se está em Fedora ou Ubuntu/RHEL
if command -v systemctl &> /dev/null; then
  systemctl daemon-reload
  systemctl enable cloudflared
  systemctl restart cloudflared
  echo "✅ Tunnel iniciado"

  sleep 2
  if systemctl is-active --quiet cloudflared; then
    echo "✅ Cloudflare tunnel ativo e funcionando"
  else
    echo "⚠️  Tunnel não parece ativo. Verifique:"
    echo "   sudo journalctl -u cloudflared -n 50"
  fi
else
  echo "⚠️  systemd não encontrado. Configure manualmente."
  echo "   Comando: sudo $CLOUDFLARE_BIN tunnel --no-autoupdate run --token '{{.Token}}'"
fi
`))

	var buf bytes.Buffer
	if err := scriptTemplate.Execute(&buf, map[string]string{"Token": p.token}); err != nil {
		return err
	}

	return os.WriteFile(setupPath, buf.Bytes(), 0755)
}

func (p *CloudflareProvisioner) printSetupInstructions() {
	log.Println("   Para configurar manualmente:")
	log.Println("   sudo /usr/bin/cloudflared tunnel --no-autoupdate run --token <SEU_TOKEN>")
}

// IsTunnelActive é público para o Agent usar no health check
func (p *CloudflareProvisioner) IsTunnelActive() bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "cloudflared")
	return cmd.Run() == nil
}

func (p *CloudflareProvisioner) DetectPackageManager() string {
	for _, pm := range []string{"dnf", "apt", "yum", "apk"} {
		if _, err := exec.LookPath(pm); err == nil {
			return pm
		}
	}
	return ""
}

func (p *CloudflareProvisioner) InstallCloudflared(pm string) error {
	log.Printf("📦 Instalando cloudflared via %s...", pm)
	switch pm {
	case "apt":
		if err := p.run("sudo", "mkdir", "-p", "--mode=0755", "/usr/share/keyrings"); err != nil {
			return err
		}
		if err := p.run("bash", "-c",
			"curl -fsSL https://pkg.cloudflare.com/cloudflare-public-v2.gpg | sudo tee /usr/share/keyrings/cloudflare-public-v2.gpg >/dev/null"); err != nil {
			return err
		}
		if err := p.run("bash", "-c",
			"echo 'deb [signed-by=/usr/share/keyrings/cloudflare-public-v2.gpg] https://pkg.cloudflare.com/cloudflared any main' | sudo tee /etc/apt/sources.list.d/cloudflared.list"); err != nil {
			return err
		}
		return p.run("sudo", "apt-get", "update", "-qq", "&&", "sudo", "apt-get", "install", "-y", "cloudflared")
	case "dnf":
		return p.run("sudo", "rpm", "-i", "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-x86_64.rpm")
	}
	return fmt.Errorf("package manager não suportado: %s", pm)
}

func (p *CloudflareProvisioner) run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *CloudflareProvisioner) RunStatus() string {
	if _, err := exec.LookPath("cloudflared"); err != nil {
		return "not_installed"
	}
	if p.isTunnelActive() {
		return "active"
	}
	if p.token != "" {
		return "configured_but_inactive"
	}
	return "not_configured"
}

func (p *CloudflareProvisioner) GetTunnelID() string {
	if p.token == "" {
		return ""
	}
	// Token JWT: eyJhIjoixx...eyJ0Ijoy... — pega o "t" claim do payload
	parts := strings.Split(p.token, ".")
	if len(parts) < 2 {
		return p.token[:16] + "..."
	}
	return p.token[len(parts[0])+1 : min(len(parts[0])+17, len(p.token))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
