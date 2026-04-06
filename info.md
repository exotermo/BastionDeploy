exodeploy/
├── api/                    # API de Controle (GO) 
│   ├── app/
│   │   ├── core/           # Configurações globais, segurança e JWT 
│   │   ├── models/         # Definição das tabelas (PostgreSQL) 
│   │   ├── schemas/        # Validação de dados (Pydantic) - A "View" dos dados 
│   │   ├── services/       # Lógica de negócio (Orquestração de deploys) 
│   │   ├── routers/        # Endpoints (deploy, status, logs) 
│   │   └── integrations/   # Clientes para LLM, Telegram e GitHub 
│   └── main.py             # Ponto de entrada da API
├── agent/                  # O Operário (Worker) 
│   ├── runners/            # Scripts de execução (Shell/Python) 
│   ├── docker/             # Gerenciamento de containers e sandboxing 
│   └── scanners/           # Auditoria (Bandit, NPM Audit) 
├── cli/                    # Exo-CLI (Typer/Go) 
│   └── main.py
├── dashboard/              # Interface Web (React/JS) 
├── infra/                  # Configurações de infraestrutura (Nginx templates) 
├── scripts/                # Scripts de instalação e utilitários de sistema
├── tests/                  # Testes automatizados
├── .env.example            # Exemplo de variáveis de ambiente
└── docker-compose.yml      # Orquestração do próprio ExoDeploy
