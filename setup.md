# BastionDeploy
Secure Self-Hosted PaaS. Kernel-level automation for the modern developer.


# Guia de Instalação e Uso do BastionDeploy

Plataforma self-hosted de automação de deploy — um "Mini-Heroku" para VPS Linux.

```
git push → webhook → API (Go) → fila Redis → Agent (Go) → git clone → docker build → docker run → Nginx/Cloudflare
```

---

## Índice

1. [Pré-requisitos](#1-pré-requisitos)
2. [Configuração rápida com exoctl (recomendado)](#2-configuração-rápida-com-exoctl-recomendado)
3. [Configuração manual](#3-configuração-manual)
4. [Infrastutura: PostgreSQL + Redis](#4-infrastrutura-postgresql--redis)
5. [Subindo a API](#5-subindo-a-api)
6. [Subindo o Agent](#6-subindo-o-agent)
7. [Subindo o Dashboard](#7-subindo-o-dashboard)
8. [Conectando o GitHub Webhook](#8-conectando-o-github-webhook)
9. [Escolhendo o modo de exposição](#9-escolhendo-o-modo-de-exposição)
10. [Comandos do exoctl](#10-comandos-do-exoctl)
11. [Teste completo (end-to-end)](#11-teste-completo-end-to-end)
12. [Estrutura de arquivos](#12-estrutura-de-arquivos)

---

## 1. Pré-requisitos

| Software | Versão mínima | Por quê |
|---|---|---|
| **Go** | 1.22+ | API e Agent são escritos em Go |
| **Git** | 2.x | Clone de repositórios pelo Agent |
| **Docker** | 24+ | Isolamento dos apps implantados |
| **PostgreSQL** | 15+ | Persistência de deploys e métricas |
| **Redis** | 7+ | Fila de jobs entre API e Agent |
| **Nginx** _(opcional)_ | 1.24+ | Proxy reverso no modo local |
| **Cloudflared** _(opcional)_ | 2024+ | Instalado automaticamente pelo Agent |

> Para desenvolvimento rápido, o `docker compose` sobe PostgreSQL + Redis automaticamente.

---

## 2. Configuração rápida com exoctl (recomendado)

```bash
# 1. Clone o repositório
git clone https://github.com/exotermo/ExoDeploy.git
cd ExoDeploy

# 2. Compile a CLI
cd cli && go build -o exoctl ./cmd/exoctl && cd ..

# 3. Rode o assistente interativo
./cli/bin/exoctl setup
```

### O que o wizard pergunta, passo a passo:

```
╔══════════════════════════════════════════════╗
║        ExoDeploy ── Setup Wizard             ║
╚══════════════════════════════════════════════╝

────── Database & Redis ──────
  DB host       [127.0.0.1]  → [Enter]
  DB port       [5432]       → [Enter]
  DB user       [exodeploy]  → [Enter]
  DB password   → sua_senha
  DB name       [exodeploy]  → [Enter]
  Redis address [127.0.0.1:6379] → [Enter]
  Redis password → sua_senha_redis

────── Discord Notifications ──────
  Cole a URL do Webhook do Discord (ou vazio para pular):
  Discord webhook URL → https://discord.com/api/webhooks/xxx

────── GitHub Webhook ──────
  Gere um secret aleatório e configure no GitHub > Settings > Webhooks
  GitHub webhook secret → [Enter = gera automático]
  Nenhum valor informado, gerando: ghsec_D4RE2g0Wtn...

────── Expose Mode ──────
  Como seus apps serão acessados?
    1) Cloudflare Tunnel (recomendado — zero portas abertas na VPS)
    2) Nginx direto na máquina (portas 80/443 abertas)
  Escolha [1/2] → 1

  ☁️  Modo: Cloudflare Tunnel
  Para obter o token:
    1. Acesse https://one.dash.cloudflare.com/
    2. Networks → Tunnels → Create a tunnel
    3. Copie o token da seção "Install and run a connector"
  Cloudflare Tunnel Token → eyJhIjo...
```

Ao final, o `exoctl` escreve os arquivos `.env` corretos em `api/.env`, `agent/.env` e `dashboard/.env`.

---

## 3. Configuração manual

Se preferir não usar a CLI, crie/editar os `.env` manualmente.

### 3.1 API (`api/.env`)

```env
PORT=8080
GIN_MODE=release

# Origens permitidas (CORS)
ALLOWED_ORIGINS=http://localhost:5173,https://seuapp.com

# Secret gerado pelo exoctl gen-secret ou manualmente
GITHUB_WEBHOOK_SECRET=um_secret_bem_grande_e_aleatorio

# PostgreSQL
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=exodeploy
DB_PASSWORD=sua_senha_segura
DB_NAME=exodeploy
DB_SSLMODE=disable

# Redis
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=sua_senha_redis
```

### 3.2 Agent (`agent/.env`)

```env
# PostgreSQL (mesmas credenciais da API)
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=exodeploy
DB_PASSWORD=sua_senha_segura
DB_NAME=exodeploy
DB_SSLMODE=disable

# Redis
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=sua_senha_redis

# Discord
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/123456/abcdef

# ─── Escolha UM dos modos abaixo ───

# Opção A: Cloudflare Tunnel (recomendado)
USE_TUNNEL_MODE=true
CLOUDFLARE_TUNNEL_TOKEN=eyJhIjoi...

# Opção B: Nginx local com Let's Encrypt
USE_TUNNEL_MODE=false
CERTBOT_EMAIL=seu@email.com
```

### 3.3 Dashboard (`dashboard/.env`)

```env
VITE_API_URL=http://127.0.0.1:8080
```

---

## 4. Infraestrutura: PostgreSQL + Redis

O `docker-compose.yml` sobe ambos os serviços:

```bash
# Copie o .env raiz se necessário — define DB_PASSWORD e REDIS_PASSWORD
# Edite .env e defina senhas reais:
#   DB_PASSWORD=sua_senha_segura
#   REDIS_PASSWORD=sua_senha_redis

docker compose up -d

# Verifique se subiram:
docker compose ps
#   exodeploy-db    → PostgreSQL :5432
#   exodeploy-redis → Redis      :6379
```

### Criar o banco e as tabelas

```bash
# Crie o banco (se não existir)
docker exec exodeploy-db createdb -U exodeploy exodeploy 2>/dev/null

# Crie a tabela de deploys
docker exec -i exodeploy-db psql -U exodeploy exodeploy <<'SQL'
CREATE TABLE IF NOT EXISTS deploys (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    app_name     VARCHAR(255) NOT NULL,
    branch       VARCHAR(255) NOT NULL DEFAULT 'main',
    commit_sha   VARCHAR(40)  NOT NULL,
    status       VARCHAR(20)  NOT NULL DEFAULT 'pending',
    triggered_by VARCHAR(255),
    created_at   TIMESTAMP    DEFAULT NOW(),
    updated_at   TIMESTAMP    DEFAULT NOW()
);
SQL
```

---

## 5. Subindo a API

```bash
cd api

# Desenvolvimento
go run cmd/api/main.go

# Produção — compila binário único
go build -o exodeploy-api ./cmd/api
./exodeploy-api
```

A API inicia na porta definida no `.env` (default `8080`).

```bash
# Teste rápido
curl http://localhost:8080
# {"service":"BastionDeploy API","version":"1.0.0"}

curl http://localhost:8080/api/v1/health
# {"status":"online","service":"ExoDeploy API","timestamp":"2025-..."}
```

---

## 6. Subindo o Agent

O Agent **deve rodar na mesma máquina onde quer implantar os aplicativos** — ele precisa de acesso ao Docker local.

```bash
cd agent

# Desenvolvimento
go run cmd/agent/main.go

# Produção
go build -o exodeploy-agent ./cmd/agent
./exodeploy-agent
```

Ao iniciar, o Agent:

1. Conecta no Redis e PostgreSQL
2. Se `USE_TUNNEL_MODE=true` — instala `cloudflared` se necessário, cria o serviço systemd e inicia o túnel
3. Entra em loop de polling na fila Redis (`BRPOP`), aguardando jobs

> **Dica:** para o Agent rodar como serviço de sistema, crie um systemd unit:

```ini
# /etc/systemd/system/exodeploy-agent.service
[Unit]
Description=ExoDeploy Agent
After=network-online.target docker.service

[Service]
Type=simple
WorkingDirectory=/opt/ExoDeploy/agent
EnvironmentFile=/opt/ExoDeploy/agent/.env
ExecStart=/opt/ExoDeploy/agent/exodeploy-agent
Restart=always

[Install]
WantedBy=multi-user.target
```

---

## 7. Subindo o Dashboard

```bash
cd dashboard

# Instala dependências
npm install

# Desenvolvimento (hot reload)
npm run dev
# → http://localhost:5173

# Produção — compila estático + serve com o Vite build
npm run build
# Output em dashboard/dist/
```

Para servir o dashboard em produção via Nginx, aponte a `root` para `dist/`.

---

## 8. Conectando o GitHub Webhook

### Passo 1: Gere um secret

```bash
./cli/bin/exoctl gen-secret
# ghsec_D4RE2g0WtnXLb33LLqCJhSWk6eMw7pQq
```

Ou use `exoctl setup` que gera automaticamente.

### Passo 2: Configure no GitHub

1. Abra seu repositório → **Settings** → **Webhooks** → **Add webhook**
2. **Payload URL**: `https://sua-api.com/api/v1/deploy/webhook`
3. **Content type**: `application/json`
4. **Secret**: o secret gerado no Passo 1
5. **Events**: marque apenas **Pushes**
6. **Add webhook**

### Passo 3: Headers opcionais

O GitHub não envia domínio ou SSL por padrão. Se quiser que o deploy configure automaticamente o domínio, configure um serviço como **ngrok** ou **webhook relay** que injeta headers customizados.

Os headers que a API entende:

| Header | Exemplo | O que faz |
|---|---|---|
| `X-Deploy-Domain` | `meuapp.exemplo.com` | Define o domínio no Nginx/Cloudflare |
| `X-Deploy-SSL` | `true` | Habilita Certificado SSL (`USE_TUNNEL_MODE=false`) |
| `X-Deploy-Port` | `3000` | Porta que o container escuta |

> **Fluxo sem domínio:** Se nenhum `X-Deploy-Domain` for enviado, o deploy funciona — o container sobe e é acessível via `http://127.0.0.1:PORTA`.

---

## 9. Escolhendo o modo de exposição

### Opção A: Cloudflare Tunnel (recomendado)

**Nenhuma porta aberta na VPS.** Todo tráfego entra via túnel outbound para Cloudflare Edge (TLS gerenciado no Edge).

```bash
# Ativar via CLI
./cli/bin/exoctl set agent USE_TUNNEL_MODE yes
./cli/bin/exoctl set agent CLOUDFLARE_TUNNEL_TOKEN "eyJh..."
```

**Obtendo o token:**

1. Acesse [https://one.dash.cloudflare.com/](https://one.dash.cloudflare.com/)
2. **Networks** → **Tunnels** → **Create a tunnel**
3. Escolha **Cloudflared** como conector
4. Copie o token da seção "Install and run a connector"

### Opção B: Nginx local + Let's Encrypt

Portas 80 e 443 abertas na VPS, certificado SSL automático via Certbot.

```bash
# Ativar via CLI
./cli/bin/exoctl set agent USE_TUNNEL_MODE no
./cli/bin/exoctl set agent CERTBOT_EMAIL "seu@email.com"
```

### Como a CLI te ajuda a trocar

```bash
# Troca de modo com um comando
./cli/bin/exoctl set agent USE_TUNNEL_MODE yes

# Verifica as configs atuais
./cli/bin/exoctl status
```

---

## 10. Comandos do exoctl

A CLI está em `cli/` e o binário em `cli/bin/exoctl` (ou `cli/exoctl` se compilado manualmente).

### `exoctl setup` — Assistente interativo

Configura tudo em sequência: banco, Redis, Discord, GitHub webhook, modo de exposição.

```bash
./cli/bin/exoctl setup
```

### `exoctl set` — Define uma configuração

```bash
# Atualiza webhook do Discord
./cli/bin/exoctl set agent DISCORD_WEBHOOK_URL https://discord.com/api/webhooks/xxx

# Ativa tunnel mode
./cli/bin/exoctl set agent USE_TUNNEL_MODE yes

# Configura domínios CORS
./cli/bin/exoctl set api ALLOWED_ORIGINS https://meuapp.com,https://admin.meuapp.com

# Troca a porta da API
./cli/bin/exoctl set api PORT 9090
```

### `exoctl get` — Consulta uma configuração

```bash
./cli/bin/exoctl get agent USE_TUNNEL_MODE
# true

./cli/bin/exoctl get agent DISCORD_WEBHOOK_URL
# http••••••••w8v1 (mascarado)
```

### `exoctl status` — Visão geral

```bash
./cli/bin/exoctl status
```

```
╔══════════════════════════════════════════════╗
║        ExoDeploy ── Status                   ║
╚══════════════════════════════════════════════╝
Install dir: /opt/ExoDeploy

── API ──
  PORT                           ✅ 8080
  GIN_MODE                       ✅ release
  ALLOWED_ORIGINS                ✅ https://meuapp.com
  GITHUB_WEBHOOK_SECRET          ✅ ghse••••••••xxxx
  DB_HOST                        ✅ 127.0.0.1
  DB_PORT                        ✅ 5432
  REDIS_ADDR                     ✅ 127.0.0.1:6379

── Agent ──
  DB_HOST                        ✅ 127.0.0.1
  REDIS_ADDR                     ✅ 127.0.0.1:6379
  DISCORD_WEBHOOK_URL            ✅ http••••••••w8v1
  USE_TUNNEL_MODE                ✅ true
  CLOUDFLARE_TUNNEL_TOKEN        ✅ eyJh••••••••xxxx
  CERTBOT_EMAIL                  ❌ não configurado
```

### `exoctl gen-secret` — Gera secret aleatório

```bash
./cli/bin/exoctl gen-secret
# ghsec_D4RE2g0WtnXLb33LLqCJhSWk6eMw7pQq
```

### `exoctl --dir` — Instalação em outro caminho

```bash
./cli/bin/exoctl status --dir=/opt/ExoDeploy
./cli/bin/exoctl setup --dir=/opt/ExoDeploy
```

---

## 11. Teste completo (end-to-end)

### 11.1 Local (simular webhook sem GitHub)

```bash
# 1. Garanta que a API e o Agent estão rodando
# Terminal 1
cd api && go run cmd/api/main.go

# Terminal 2
cd agent && go run cmd/agent/main.go

# Terminal 3 — simula um push
curl -X POST http://localhost:8080/api/v1/deploy/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -H "X-Hub-Signature-256: sha256=$(echo -n '{"ref":"refs/heads/main","head_commit":{"id":"abc123"},"repository":{"name":"minha-app","clone_url":"https://github.com/user/minha-app.git"},"pusher":{"name":"test"}}' | openssl dgst -sha256 -hmac 'SEU_SECRET' -hex | cut -d' ' -f2 | sed 's/^.*=//' | xxd -r -p | openssl dgst -sha256 -hmac 'SEU_SECRET' -hex | awk '{print "sha256="$2}')" \
  -d '{
    "ref": "refs/heads/main",
    "head_commit": {"id": "abc123"},
    "repository": {"name": "minha-app", "clone_url": "https://github.com/user/minha-app.git"},
    "pusher": {"name": "dev"}
  }'
```

### 11.2 Verificar o deploy

```bash
# Histórico de deploys
curl http://localhost:8080/api/v1/deploys

# Status de um app específico
curl http://localhost:8080/api/v1/deploy/status/minha-app

# Stats gerais
curl http://localhost:8080/api/v1/stats
```

### 11.3 Fluxo real com GitHub

```
1. Faça push para qualquer repo configurado
        ↓
2. GitHub envia webhook → API valida HMAC
        ↓
3. API enfileira job no Redis
        ↓
4. Agent recebe job → clone → build → run → Nginx
        ↓
5. Agent atualiza status no DB
        ↓
6. Notificação no Discord
```

---

## 12. Estrutura de arquivos

```
ExoDeploy/
├── api/                          # Go — API de controle
│   ├── cmd/api/main.go           # Entry point
│   ├── internal/
│   │   ├── domain/deploy.go      # Domain models + interfaces
│   │   ├── handler/              # HTTP handlers (webhook, stats, health)
│   │   └── repository/           # PostgreSQL repository
│   ├── pkg/
│   │   ├── config/config.go      # Env config loader
│   │   └── middleware/           # HMAC, CORS, security headers
│   └── .env
│
├── agent/                        # Go — Agent de deploy
│   ├── cmd/agent/main.go         # Entry point
│   ├── internal/
│   │   ├── worker/               # Redis queue consumer
│   │   ├── runner/               # Docker executor (git clone, build, run)
│   │   ├── provisioner/          # Nginx vhost + Cloudflare Tunnel
│   │   └── notifier/             # Discord webhook notifier
│   └── .env
│
├── cli/                          # Go — exoctl CLI
│   ├── cmd/exoctl/main.go        # CLI entry point + commands
│   └── internal/
│       ├── config/env.go         # .env read/write
│       └── wizard/wizard.go      # Interactive setup wizard
│
├── dashboard/                    # React/Vite/TS — Painel web
│   └── src/
│       ├── components/           # UI components
│       └── hooks/                # API data fetching
│
├── docker-compose.yml            # PostgreSQL + Redis
├── setup.md                      # Este arquivo
└── .env                          # Senhas da infra (DB_PASSWORD, REDIS_PASSWORD)
```

---

## Problemas comuns

| Sintoma | Causa | Solução |
|---|---|---|
| `Erro: não foi possível encontrar o ExoDeploy` | Executando `exoctl` de pasta errada | Use `--dir=/caminho/do/repo` |
| `assinatura inválida` | Secret do `.env` ≠ secret configurado no GitHub | `exoctl set api GITHUB_WEBHOOK_SECRET <novo-secret>` |
| Agent não conecta no Redis | Senha errada ou serviço parado | `docker compose ps` + verifique `REDIS_PASSWORD` |
| `nginx -t` falha | Nginx não instalado | `dnf install nginx` ou `apt install nginx` |
| Certbot falha | Porta 80 não acessível | Verifique firewall e DNS do domínio |
| Container não sobe | Porta já em uso | Mude a porta com `X-Deploy-Port` ou pare o container conflituoso |
