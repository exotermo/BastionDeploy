package runner

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type DockerRunner struct{}

func NewDockerRunner() *DockerRunner {
	return &DockerRunner{}
}

// Deploy clona o repo, builda e sobe o container
func (r *DockerRunner) Deploy(ctx context.Context, appName, repoURL, branch string) error {
	workDir := filepath.Join("/tmp/bastiondeploy", appName)

	// 1. Limpa diretório anterior se existir
	os.RemoveAll(workDir)

	// 2. Clona o repositório
	log.Printf("📥 Clonando %s (branch: %s)", repoURL, branch)
	if err := run(ctx, "git", "clone", "--depth=1", "--branch", branch, repoURL, workDir); err != nil {
		return fmt.Errorf("git clone falhou: %w", err)
	}

	// 3. Gera Dockerfile automático se não existir
	dockerfilePath := filepath.Join(workDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		log.Printf("📝 Dockerfile não encontrado, gerando automaticamente para %s", appName)
		if err := generateDockerfile(workDir, appName); err != nil {
			return fmt.Errorf("erro ao gerar Dockerfile: %w", err)
		}
	}

	// 4. Build da imagem
	imageName := fmt.Sprintf("bastiondeploy-%s:latest", appName)
	log.Printf("🔨 Buildando imagem %s", imageName)
	if err := run(ctx, "docker", "build", "-t", imageName, workDir); err != nil {
		return fmt.Errorf("docker build falhou: %w", err)
	}

	// 5. Para container antigo se existir
	containerName := fmt.Sprintf("app-%s", appName)
	run(ctx, "docker", "stop", containerName)
	run(ctx, "docker", "rm", containerName)

	// 6. Sobe o novo container
	log.Printf("▶️  Subindo container %s", containerName)
	if err := run(ctx, "docker", "run", "-d",
		"--name", containerName,
		"--restart", "unless-stopped",
		imageName,
	); err != nil {
		return fmt.Errorf("docker run falhou: %w", err)
	}

	return nil
}

// generateDockerfile detecta a linguagem e gera um Dockerfile básico
func generateDockerfile(workDir, appName string) error {
	var content string

	switch {
	case fileExists(filepath.Join(workDir, "package.json")):
		content = `FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build 2>/dev/null || true
CMD ["node", "index.js"]`

	case fileExists(filepath.Join(workDir, "requirements.txt")):
		content = `FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
CMD ["python", "main.py"]`

	case fileExists(filepath.Join(workDir, "go.mod")):
		content = `FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/...

FROM alpine:latest
COPY --from=builder /app/app .
CMD ["./app"]`

	default:
		content = `FROM alpine:latest
COPY . /app
WORKDIR /app
CMD ["sh", "-c", "echo 'App running'"]`
	}

	return os.WriteFile(filepath.Join(workDir, "Dockerfile"), []byte(content), 0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// run executa um comando e loga o output
func run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}