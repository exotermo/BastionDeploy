#!/bin/bash

# Caminho base para a API
API_DIR="."

echo "🚀 Estruturando a API ExoDeploy em Go (Padrão Enterprise)..."

# 1. Criando os pontos de entrada (onde o binário nasce)
mkdir -p "$API_DIR/cmd/api"

# 2. Criando o coração do sistema (código privado que não pode ser importado por outros)
# Aqui fica a regra de negócio do ExoDeploy
mkdir -p "$API_DIR/internal/domain"      # Entidades/Interfaces (Ex: Deploy, Server, User)
mkdir -p "$API_DIR/internal/service"     # Lógica: Onde a IA e o Docker SDK vão agir
mkdir -p "$API_DIR/internal/handler"     # Seus "Routers" e Handlers HTTP (FastAPI style)
mkdir -p "$API_DIR/internal/repository"  # Conexão com Banco de Dados (PostgreSQL)
mkdir -p "$API_DIR/internal/config"      # Leitura de .env e configurações globais

# 3. Criando pacotes utilitários (código que pode ser compartilhado)
mkdir -p "$API_DIR/pkg/logger"           # Logs estruturados para o dashboard
mkdir -p "$API_DIR/pkg/auth"             # Lógica de JWT e RBAC

# 4. Inicializando arquivos básicos
touch "$API_DIR/cmd/api/main.go"
touch "$API_DIR/go.mod"

echo "✅ Estrutura Go criada!"
echo "Sugestão: Rode 'go mod init github.com/exotermo/exodeploy/api' agora."
