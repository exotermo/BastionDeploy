#!/bin/bash
# ===========================================
# BastionDeploy — Dashboard Setup Script
# ===========================================

set -e  # Para o script se qualquer comando falhar

DASHBOARD_DIR="/home/xande/Documentos/ProjetosFedora/ExoDeploy/dashboard"

echo "🚀 Criando estrutura do dashboard..."

# ===========================================
# 1. ESTRUTURA DE PASTAS
# ===========================================
mkdir -p $DASHBOARD_DIR/src/components
mkdir -p $DASHBOARD_DIR/src/styles
mkdir -p $DASHBOARD_DIR/nginx

echo "✅ Pastas criadas"

# ===========================================
# 2. package.json
# ===========================================
cat > $DASHBOARD_DIR/package.json << 'EOF'
{
  "name": "bastiondeploy-dashboard",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "lucide-react": "^0.487.0",
    "recharts": "^2.15.2",
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  },
  "devDependencies": {
    "@types/react": "^18.3.1",
    "@types/react-dom": "^18.3.1",
    "@vitejs/plugin-react": "^4.7.0",
    "@tailwindcss/vite": "^4.1.12",
    "tailwindcss": "^4.1.12",
    "typescript": "^5.8.3",
    "vite": "^6.3.5"
  }
}
EOF

echo "✅ package.json criado"

# ===========================================
# 3. vite.config.ts
# ===========================================
cat > $DASHBOARD_DIR/vite.config.ts << 'EOF'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
EOF

# ===========================================
# 4. tsconfig.json
# ===========================================
cat > $DASHBOARD_DIR/tsconfig.json << 'EOF'
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src"]
}
EOF

echo "✅ Configs de build criadas"

# ===========================================
# 5. index.html
# ===========================================
cat > $DASHBOARD_DIR/index.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>BastionDeploy — Dashboard</title>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet" />
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
EOF

# ===========================================
# 6. src/styles/index.css
# ===========================================
cat > $DASHBOARD_DIR/src/styles/index.css << 'EOF'
@import "tailwindcss";

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  background-color: #0f1117;
  color: #ffffff;
  font-family: 'Inter', sans-serif;
  -webkit-font-smoothing: antialiased;
}

/* Scrollbar personalizada */
::-webkit-scrollbar {
  width: 6px;
}
::-webkit-scrollbar-track {
  background: #0a0c10;
}
::-webkit-scrollbar-thumb {
  background: #ffffff20;
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: #ffffff40;
}
EOF

# ===========================================
# 7. src/main.tsx
# ===========================================
cat > $DASHBOARD_DIR/src/main.tsx << 'EOF'
import { createRoot } from 'react-dom/client'
import App from './App'
import './styles/index.css'

createRoot(document.getElementById('root')!).render(<App />)
EOF

# ===========================================
# 8. src/App.tsx
# ===========================================
cat > $DASHBOARD_DIR/src/App.tsx << 'EOF'
import { Sidebar } from './components/Sidebar'
import { StatsCards } from './components/StatsCards'
import { DeploymentsTable } from './components/DeploymentsTable'
import { DeployChart } from './components/DeployChart'
import { AppStatusPanel } from './components/AppStatusPanel'
import { AIAssistant } from './components/AIAssistant'

export default function App() {
  return (
    <div className="flex h-screen bg-[#0f1117] overflow-hidden">
      <Sidebar />
      <div className="flex-1 overflow-y-auto">
        <div className="p-8 max-w-[1800px]">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">Dashboard</h1>
            <p className="text-gray-400">Monitor your deployments and apps in real-time</p>
          </div>
          <div className="mb-8">
            <StatsCards />
          </div>
          <div className="mb-8">
            <AIAssistant />
          </div>
          <div className="mb-8">
            <DeploymentsTable />
          </div>
          <div className="grid grid-cols-2 gap-6">
            <DeployChart />
            <AppStatusPanel />
          </div>
        </div>
      </div>
    </div>
  )
}
EOF

# ===========================================
# 9. COMPONENTES
# ===========================================

cat > $DASHBOARD_DIR/src/components/Sidebar.tsx << 'EOF'
import { Rocket, LayoutDashboard, GitBranch, Package, ScrollText, Settings, User } from 'lucide-react'

export function Sidebar() {
  const navItems = [
    { icon: LayoutDashboard, label: 'Dashboard', active: true },
    { icon: GitBranch, label: 'Deployments', active: false },
    { icon: Package, label: 'Apps', active: false },
    { icon: ScrollText, label: 'Logs', active: false },
    { icon: Settings, label: 'Settings', active: false },
  ]

  return (
    <div className="w-64 bg-[#0a0c10] border-r border-[#ffffff10] flex flex-col h-screen">
      <div className="p-6 border-b border-[#ffffff10]">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-gradient-to-br from-[#00d4ff] to-[#0088ff] rounded-lg flex items-center justify-center">
            <Rocket className="w-5 h-5 text-white" />
          </div>
          <span className="text-xl font-semibold text-white">BastionDeploy</span>
        </div>
      </div>
      <nav className="flex-1 p-4">
        <ul className="space-y-1">
          {navItems.map((item) => (
            <li key={item.label}>
              <button
                className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
                  item.active
                    ? 'bg-[#00d4ff15] text-[#00d4ff] border border-[#00d4ff30]'
                    : 'text-gray-400 hover:text-white hover:bg-[#ffffff08]'
                }`}
              >
                <item.icon className="w-5 h-5" />
                <span className="font-medium">{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>
      <div className="p-4 border-t border-[#ffffff10]">
        <div className="flex items-center gap-3 px-4 py-3 rounded-lg bg-[#ffffff05] border border-[#ffffff08]">
          <div className="w-10 h-10 bg-gradient-to-br from-[#00d4ff] to-[#0088ff] rounded-full flex items-center justify-center">
            <User className="w-5 h-5 text-white" />
          </div>
          <div>
            <p className="text-white font-medium text-sm">xande</p>
            <p className="text-gray-500 text-xs">admin@bastiondeploy</p>
          </div>
        </div>
      </div>
    </div>
  )
}
EOF

cat > $DASHBOARD_DIR/src/components/StatsCards.tsx << 'EOF'
import { TrendingUp, CheckCircle2, Package, Clock } from 'lucide-react'

export function StatsCards() {
  const stats = [
    { icon: TrendingUp,   label: 'Total Deploys', value: '42',       color: '#00d4ff' },
    { icon: CheckCircle2, label: 'Success Rate',  value: '94%',      color: '#00ff88' },
    { icon: Package,      label: 'Active Apps',   value: '7',        color: '#00d4ff' },
    { icon: Clock,        label: 'Last Deploy',   value: '2 min ago',color: '#ffaa00' },
  ]

  return (
    <div className="grid grid-cols-4 gap-6">
      {stats.map((stat) => (
        <div
          key={stat.label}
          className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl p-6 hover:border-[#ffffff20] transition-all"
          style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}
        >
          <div className="flex items-start justify-between">
            <div>
              <p className="text-gray-400 text-sm mb-2">{stat.label}</p>
              <p className="text-3xl font-bold text-white">{stat.value}</p>
            </div>
            <div className="p-3 rounded-lg" style={{ backgroundColor: `${stat.color}15` }}>
              <stat.icon className="w-6 h-6" style={{ color: stat.color }} />
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
EOF

cat > $DASHBOARD_DIR/src/components/AIAssistant.tsx << 'EOF'
import { Bot, AlertCircle } from 'lucide-react'

export function AIAssistant() {
  return (
    <div
      className="bg-gradient-to-br from-[#0a0c10] to-[#0f1117] border border-[#00d4ff30] rounded-xl p-6 relative overflow-hidden"
      style={{ boxShadow: '0 0 40px rgba(0,212,255,0.15), 0 4px 24px rgba(0,0,0,0.4)' }}
    >
      <div className="absolute -top-20 -right-20 w-40 h-40 bg-[#00d4ff] rounded-full blur-[100px] opacity-20" />
      <div className="relative z-10">
        <div className="flex items-start gap-4">
          <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-[#00d4ff] to-[#0088ff] rounded-lg flex items-center justify-center shadow-lg shadow-[#00d4ff30]">
            <Bot className="w-6 h-6 text-white" />
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-3">
              <h3 className="text-white font-semibold">AI Agent</h3>
              <span className="text-xs text-[#00d4ff] bg-[#00d4ff15] px-2 py-0.5 rounded-full border border-[#00d4ff30]">
                Active
              </span>
            </div>
            <div className="bg-[#0f1117] border border-[#ffffff10] rounded-lg p-4 text-sm" style={{ fontFamily: 'JetBrains Mono, monospace' }}>
              <div className="flex items-start gap-2 mb-2">
                <AlertCircle className="w-4 h-4 text-[#ff4444] flex-shrink-0 mt-0.5" />
                <div className="text-gray-300">
                  <span className="text-[#ff4444] font-medium">Deploy failed</span> on{' '}
                  <span className="text-[#00d4ff]">meu-bot</span>
                </div>
              </div>
              <div className="pl-6 space-y-2">
                <p className="text-gray-400">
                  Missing environment variable{' '}
                  <code className="text-[#00ff88] bg-[#00ff8815] px-1.5 py-0.5 rounded">DATABASE_URL</code>
                </p>
                <div className="border-l-2 border-[#00d4ff30] pl-3 mt-3">
                  <p className="text-[#00d4ff] text-xs font-medium mb-1">💡 Suggestion:</p>
                  <p className="text-gray-300">
                    Add it to your <code className="text-gray-400">.env</code> file and redeploy.
                  </p>
                </div>
              </div>
            </div>
            <div className="mt-3 text-xs text-gray-500">Last analysis: 2 minutes ago</div>
          </div>
        </div>
      </div>
    </div>
  )
}
EOF

cat > $DASHBOARD_DIR/src/components/DeploymentsTable.tsx << 'EOF'
import { GitBranch, User, Clock } from 'lucide-react'

interface Deployment {
  app: string
  branch: string
  commit: string
  status: 'SUCCESS' | 'FAILED' | 'PENDING' | 'RUNNING'
  triggeredBy: string
  time: string
}

const deployments: Deployment[] = [
  { app: 'meu-bot',      branch: 'main', commit: 'abc123', status: 'SUCCESS', triggeredBy: 'xande', time: '2 min ago'  },
  { app: 'api-finance',  branch: 'dev',  commit: 'def456', status: 'FAILED',  triggeredBy: 'xande', time: '1 hour ago' },
  { app: 'dashboard',    branch: 'main', commit: 'ghi789', status: 'PENDING', triggeredBy: 'xande', time: 'just now'   },
  { app: 'telegram-bot', branch: 'main', commit: 'jkl012', status: 'RUNNING', triggeredBy: 'xande', time: '30s ago'    },
]

const statusStyle: Record<Deployment['status'], string> = {
  SUCCESS: 'bg-[#00ff8815] text-[#00ff88] border-[#00ff8830]',
  FAILED:  'bg-[#ff444415] text-[#ff4444] border-[#ff444430]',
  PENDING: 'bg-[#ffaa0015] text-[#ffaa00] border-[#ffaa0030]',
  RUNNING: 'bg-[#00d4ff15] text-[#00d4ff] border-[#00d4ff30] animate-pulse',
}

export function DeploymentsTable() {
  return (
    <div className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl overflow-hidden"
         style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}>
      <div className="p-6 border-b border-[#ffffff10]">
        <h2 className="text-xl font-semibold text-white">Recent Deployments</h2>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[#ffffff10]">
              {['App','Branch','Commit','Status','Triggered by','Time'].map(h => (
                <th key={h} className="text-left px-6 py-4 text-sm font-medium text-gray-400">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {deployments.map((d, i) => (
              <tr key={i} className="border-b border-[#ffffff08] hover:bg-[#ffffff05] transition-colors">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2">
                    <div className="w-8 h-8 bg-[#00d4ff15] rounded-lg flex items-center justify-center">
                      <GitBranch className="w-4 h-4 text-[#00d4ff]" />
                    </div>
                    <span className="text-white font-medium">{d.app}</span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span className="text-gray-300" style={{ fontFamily: 'JetBrains Mono, monospace' }}>{d.branch}</span>
                </td>
                <td className="px-6 py-4">
                  <code className="text-gray-400 text-sm bg-[#ffffff08] px-2 py-1 rounded" style={{ fontFamily: 'JetBrains Mono, monospace' }}>{d.commit}</code>
                </td>
                <td className="px-6 py-4">
                  <span className={`px-3 py-1 rounded-full text-xs font-medium border ${statusStyle[d.status]}`}>{d.status}</span>
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2 text-gray-300">
                    <User className="w-4 h-4" />{d.triggeredBy}
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2 text-gray-400">
                    <Clock className="w-4 h-4" />{d.time}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
EOF

cat > $DASHBOARD_DIR/src/components/DeployChart.tsx << 'EOF'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

const data = [
  { day: 'Mon', deploys: 5 },
  { day: 'Tue', deploys: 8 },
  { day: 'Wed', deploys: 6 },
  { day: 'Thu', deploys: 10 },
  { day: 'Fri', deploys: 7 },
  { day: 'Sat', deploys: 3 },
  { day: 'Sun', deploys: 4 },
]

export function DeployChart() {
  return (
    <div className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl p-6"
         style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}>
      <h3 className="text-lg font-semibold text-white mb-6">Deploys Last 7 Days</h3>
      <ResponsiveContainer width="100%" height={240}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#ffffff10" />
          <XAxis dataKey="day" stroke="#6b7280" style={{ fontSize: '12px' }} />
          <YAxis stroke="#6b7280" style={{ fontSize: '12px' }} />
          <Tooltip
            contentStyle={{ backgroundColor: '#0a0c10', border: '1px solid #ffffff20', borderRadius: '8px', color: '#fff' }}
            cursor={{ fill: '#ffffff08' }}
          />
          <Bar dataKey="deploys" fill="#00d4ff" radius={[8, 8, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
EOF

cat > $DASHBOARD_DIR/src/components/AppStatusPanel.tsx << 'EOF'
import { Circle } from 'lucide-react'

interface AppStatus { name: string; status: 'UP' | 'DOWN'; uptime: string }

const apps: AppStatus[] = [
  { name: 'meu-bot',         status: 'UP',   uptime: '99.9%' },
  { name: 'api-finance',     status: 'DOWN', uptime: '0%'    },
  { name: 'dashboard',       status: 'UP',   uptime: '99.2%' },
  { name: 'telegram-bot',    status: 'UP',   uptime: '100%'  },
  { name: 'analytics-api',   status: 'UP',   uptime: '98.7%' },
  { name: 'webhook-service', status: 'UP',   uptime: '99.5%' },
  { name: 'cron-worker',     status: 'UP',   uptime: '100%'  },
]

export function AppStatusPanel() {
  return (
    <div className="bg-[#0a0c10] border border-[#ffffff10] rounded-xl p-6"
         style={{ boxShadow: '0 4px 24px rgba(0,0,0,0.4)' }}>
      <h3 className="text-lg font-semibold text-white mb-6">App Status</h3>
      <div className="space-y-4">
        {apps.map((app) => (
          <div key={app.name}
               className="flex items-center justify-between p-3 rounded-lg bg-[#ffffff05] border border-[#ffffff08] hover:border-[#ffffff15] transition-all">
            <div className="flex items-center gap-3">
              <Circle className={`w-3 h-3 ${app.status === 'UP' ? 'fill-[#00ff88] text-[#00ff88]' : 'fill-[#ff4444] text-[#ff4444]'}`} />
              <span className="text-white font-medium" style={{ fontFamily: 'JetBrains Mono, monospace' }}>{app.name}</span>
            </div>
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-400">Uptime: {app.uptime}</span>
              <span className={`px-2 py-1 rounded text-xs font-medium ${
                app.status === 'UP'
                  ? 'bg-[#00ff8815] text-[#00ff88] border border-[#00ff8830]'
                  : 'bg-[#ff444415] text-[#ff4444] border border-[#ff444430]'
              }`}>{app.status}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
EOF

echo "✅ Componentes criados"

# ===========================================
# 10. Dockerfile
# ===========================================
cat > $DASHBOARD_DIR/Dockerfile << 'EOF'
# Stage 1 — build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2 — serve com nginx
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
EOF

# ===========================================
# 11. nginx/nginx.conf
# ===========================================
cat > $DASHBOARD_DIR/nginx/nginx.conf << 'EOF'
server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    # Headers de segurança
    add_header X-Frame-Options "DENY";
    add_header X-Content-Type-Options "nosniff";
    add_header Referrer-Policy "strict-origin-when-cross-origin";
    add_header X-XSS-Protection "1; mode=block";

    # SPA — redireciona tudo pro index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache para assets estáticos
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Bloqueia acesso a arquivos sensíveis
    location ~ /\. {
        deny all;
    }
}
EOF

# ===========================================
# 12. .dockerignore
# ===========================================
cat > $DASHBOARD_DIR/.dockerignore << 'EOF'
node_modules
dist
.env
*.log
EOF

echo "✅ Docker e Nginx configurados"
echo ""
echo "🎉 Setup completo! Agora rode:"
echo ""
echo "  cd $DASHBOARD_DIR"
echo "  npm install"
echo "  npm run dev"
echo ""
