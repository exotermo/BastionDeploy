#!/bin/bash

# Definindo o diretório base
BASE_DIR="/home/xande/Documentos/ProjetosFedora/ExoDeploy"

echo "Iniciando a criação da estrutura ExoDeploy em: $BASE_DIR"

# Criando a estrutura de pastas da API (Cérebro)
mkdir -p "$BASE_DIR/api/app/core"
mkdir -p "$BASE_DIR/api/app/models"
mkdir -p "$BASE_DIR/api/app/schemas"
mkdir -p "$BASE_DIR/api/app/services"
mkdir -p "$BASE_DIR/api/app/routers"
mkdir -p "$BASE_DIR/api/app/integrations"

# Criando a estrutura do Agent (Operário)
mkdir -p "$BASE_DIR/agent/runners"
mkdir -p "$BASE_DIR/agent/docker"
mkdir -p "$BASE_DIR/agent/scanners"

# Criando pastas de suporte e ecossistema
mkdir -p "$BASE_DIR/cli"
mkdir -p "$BASE_DIR/dashboard"
mkdir -p "$BASE_DIR/infra"
mkdir -p "$BASE_DIR/scripts"
mkdir -p "$BASE_DIR/tests"

# Criando arquivos de configuração de infraestrutura
touch "$BASE_DIR/.env.example"
touch "$BASE_DIR/docker-compose.yml"

echo "Estrutura criada com sucesso!"
ls -R "$BASE_DIR"
