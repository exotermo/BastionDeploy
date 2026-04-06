package provisioner

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

const (
	cloudflaredService = "/etc/systemd/system/cloudflared.service"
	cloudflaredEnvFile = "/etc/default/cloudflared"
)

// CloudflareProvisioner gerencia o cloudflared tunnel na VPS
type CloudflareProvisioner struct {
	token string
}

func NewCloudflareProvisioner(token string) *CloudflareProvisioner {
	return &CloudflareProvisioner{token: token}
}

// EnsureCloudflared verifica se cloudflared está instalado, senão instala
func (p *CloudflareProvisioner) EnsureCloudflared() error {
	if _, err := exec.LookPath("cloudflared"); err == nil {
		log.Println("☁️  cloudflared já está instalado")
		return nil
	}

	log.Println("📦 Instalando cloudflared...")
	return p.installCloudflared()
}

// SetupTunnel configura e inicia o cloudflared tunnel com o token fornecido
func (p *CloudflareProvisioner) SetupTunnel() error {
	if err := p.EnsureCloudflared(); err != nil {
		return fmt.Errorf("erro ao instalar cloudflared: %w", err)
	}

	// Verifica se o túnel já está configurado
	if _, err := os.Stat(cloudflaredEnvFile); err == nil {
		content, _ := os.ReadFile(cloudflaredEnvFile)
		if bytes.Contains(content, []byte("TUNNEL_TOKEN")) {
			log.Println("☁️  Cloudflare tunnel já está configurado, pulando")
			return nil
		}
	}

	if p.token == "" {
		return fmt.Errorf("CLOUDFLARE_TUNNEL_TOKEN não definido — obrigatório para tunnel mode")
	}

	// Cria o systemd service para o tunnel
	if err := p.writeSystemdService(); err != nil {
		return fmt.Errorf("erro ao criar systemd service: %w", err)
	}

	// Grava o token de autenticação
	if err := os.WriteFile(cloudflaredEnvFile,
		[]byte(fmt.Sprintf("TUNNEL_TOKEN=%s\n", p.token)),
		0600,
	); err != nil {
		return fmt.Errorf("erro ao gravar token: %w", err)
	}

	// Reload + enable + start
	if err := p.runServiceCommands("daemon-reload", "enable", "start"); err != nil {
		return fmt.Errorf("erro ao iniciar tunnel: %w", err)
	}

	log.Println("✅ Cloudflare tunnel ativo")
	return nil
}

// IsTunnelActive verifica se o serviço cloudflared está rodando
func (p *CloudflareProvisioner) IsTunnelActive() bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "cloudflared")
	return cmd.Run() == nil
}

// --- internals ---

func (p *CloudflareProvisioner) installCloudflared() error {
	// Detecta o gerenciador de pacotes e instala
	if _, err := exec.LookPath("dnf"); err == nil {
		return p.runInstall("dnf", "install", "-y", "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-x86_64.rpm")
	}
	if _, err := exec.LookPath("apt"); err == nil {
		// curl + dpkg para .deb
		if err := p.run("curl", "-fSL", "-o", "/tmp/cloudflared.deb",
			"https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb"); err != nil {
			return err
		}
		defer os.Remove("/tmp/cloudflared.deb")
		return p.run("dpkg", "-i", "/tmp/cloudflared.deb")
	}

	return fmt.Errorf("gerenciador de pacotes não suportado — instale cloudflared manualmente")
}

func (p *CloudflareProvisioner) runInstall(name string, args ...string) error {
	return p.run(name, args...)
}

var cloudflaredServiceTemplate = template.Must(template.New("cloudflared").Parse(`[Unit]
Description=Cloudflare Tunnel
After=network-online.target

[Service]
Type=simple
EnvironmentFile=/etc/default/cloudflared
ExecStart=/usr/local/bin/cloudflared tunnel --no-autoupdate run
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
`))

func (p *CloudflareProvisioner) writeSystemdService() error {
	var buf bytes.Buffer
	if err := cloudflaredServiceTemplate.Execute(&buf, nil); err != nil {
		return err
	}
	return os.WriteFile(cloudflaredService, buf.Bytes(), 0644)
}

func (p *CloudflareProvisioner) runServiceCommands(subcommands ...string) error {
	for _, sub := range subcommands {
		if err := p.run("systemctl", sub, "cloudflared"); err != nil {
			return fmt.Errorf("systemctl %s falhou: %w", sub, err)
		}
	}
	return nil
}

func (p *CloudflareProvisioner) run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
