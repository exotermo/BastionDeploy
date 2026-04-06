import { useState, useEffect, useRef } from 'react'
import { useDeploys, Deploy } from '../hooks/useDeploys'

type PipelineStage = {
  id: string
  label: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'skipped'
  startedAt?: number
  finishedAt?: number
}

type DeployLog = {
  stage: string
  message: string
  timestamp: string
  level: 'info' | 'warn' | 'error' | 'success'
}

const PIPELINE_STAGES: { id: string; label: string }[] = [
  { id: 'webhook', label: 'Webhook' },
  { id: 'clone', label: 'Git Clone' },
  { id: 'build', label: 'Docker Build' },
  { id: 'deploy', label: 'Docker Run' },
  { id: 'nginx', label: 'Nginx' },
  { id: 'notify', label: 'Notify' },
]

function generatePipelineLogs(deploy: Deploy): PipelineStage[] {
  const stages: PipelineStage[] = PIPELINE_STAGES.map((s) => ({
    ...s,
    status: 'pending' as const,
  }))

  switch (deploy.status) {
    case 'success':
      stages.forEach((s, i) => {
        s.status = 'success'
        s.startedAt = i * 2000
        s.finishedAt = (i + 1) * 2000
      })
      break
    case 'failed':
      stages.forEach((s, i) => {
        if (i < 3) {
          s.status = 'success'
          s.startedAt = i * 2000
          s.finishedAt = (i + 1) * 2000
        } else if (i === 3) {
          s.status = 'failed'
          s.startedAt = i * 2000
          s.finishedAt = (i + 1) * 2000
        } else {
          s.status = 'skipped'
        }
      })
      break
    case 'running':
      stages.forEach((s, i) => {
        if (i < 2) {
          s.status = 'success'
          s.startedAt = i * 2000
          s.finishedAt = (i + 1) * 2000
        } else if (i === 2) {
          s.status = 'running'
          s.startedAt = Date.now()
        } else {
          s.status = 'pending'
        }
      })
      break
    case 'pending':
    default:
      break
  }

  return stages
}

const DEPLOY_LOG_TEMPLATES: Record<string, DeployLog[]> = {
  webhook: [
    { stage: 'webhook', message: 'Recebendo webhook do GitHub...', timestamp: '', level: 'info' },
    { stage: 'webhook', message: 'Validando assinatura HMAC-SHA256', timestamp: '', level: 'info' },
    { stage: 'webhook', message: 'Assinatura verificada com sucesso', timestamp: '', level: 'success' },
  ],
  clone: [
    { stage: 'clone', message: 'Iniciando git clone --depth=1 --branch main', timestamp: '', level: 'info' },
    { stage: 'clone', message: 'Repositório clonado com sucesso', timestamp: '', level: 'success' },
    { stage: 'clone', message: 'Detectando linguagem do projeto...', timestamp: '', level: 'info' },
  ],
  build: [
    { stage: 'build', message: 'Executando docker build...', timestamp: '', level: 'info' },
    { stage: 'build', message: 'Resolvendo dependências do projeto', timestamp: '', level: 'info' },
    { stage: 'build', message: 'Imagem construída com sucesso', timestamp: '', level: 'success' },
  ],
  deploy: [
    { stage: 'deploy', message: 'Parando container antigo (se existir)', timestamp: '', level: 'info' },
    { stage: 'deploy', message: 'Iniciando container com docker run', timestamp: '', level: 'info' },
    { stage: 'deploy', message: 'Container ativo e respondendo', timestamp: '', level: 'success' },
  ],
  nginx: [
    { stage: 'nginx', message: 'Gerando vhost do Nginx...', timestamp: '', level: 'info' },
    { stage: 'nginx', message: 'nginx -t OK', timestamp: '', level: 'success' },
    { stage: 'nginx', message: 'Nginx reload executado', timestamp: '', level: 'info' },
    { stage: 'nginx', message: 'Cloudflare Tunnel: rota configurada', timestamp: '', level: 'info' },
  ],
  notify: [
    { stage: 'notify', message: 'Enviando notificação para o Discord', timestamp: '', level: 'info' },
    { stage: 'notify', message: 'Deploy concluído com sucesso', timestamp: '', level: 'success' },
  ],
}

const ERROR_LOG_TEMPLATES: DeployLog[] = [
  { stage: 'deploy', message: 'Erro: container falhou ao iniciar', timestamp: '', level: 'error' },
  { stage: 'deploy', message: 'Verifique os logs do container para mais detalhes', timestamp: '', level: 'warn' },
]

function buildFullLogs(deploy: Deploy): DeployLog[] {
  const stages = generatePipelineLogs(deploy)
  let logs: DeployLog[] = []
  let t = 0

  for (const stage of stages) {
    const templates =
      stage.status === 'failed' && stage.id === 'deploy'
        ? ERROR_LOG_TEMPLATES
        : DEPLOY_LOG_TEMPLATES[stage.id] ?? []

    for (const log of templates) {
      t += Math.floor(Math.random() * 1000) + 500
      const date = new Date(new Date(deploy.created_at).getTime() + t)
      logs.push({
        ...log,
        timestamp: date.toLocaleTimeString('pt-BR', {
          hour: '2-digit',
          minute: '2-digit',
          second: '2-digit',
        }),
      })
    }
  }

  if (deploy.status === 'success') {
    logs.push({
      stage: 'notify',
      message: `Deploy ${deploy.app_name} (${deploy.commit_sha.slice(0, 7)}) finalizado com sucesso`,
      timestamp: '',
      level: 'success',
    })
    logs[logs.length - 1].timestamp = (() => {
      t += 300
      const date = new Date(new Date(deploy.created_at).getTime() + t)
      return date.toLocaleTimeString('pt-BR', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      })
    })()
  }

  return logs
}

interface DeployLogProps {
  deploy: Deploy
}

export function DeployLog({ deploy }: DeployLogProps) {
  const [logs, setLogs] = useState<DeployLog[]>([])
  const logEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setLogs(buildFullLogs(deploy))
  }, [deploy])

  useEffect(() => {
    if (deploy.status === 'running') {
      const interval = setInterval(() => {
        const base = buildFullLogs(deploy)
        setLogs((prev) => {
          if (prev.length >= base.length) return prev
          // Simula logs aparecendo em tempo real
          return base.slice(0, Math.min(base.length, prev.length + 2))
        })
      }, 1500)
      return () => clearInterval(interval)
    }
  }, [deploy])

  useEffect(() => {
    logEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [logs.length])

  const stages = generatePipelineLogs(deploy)

  const levelDot: Record<string, string> = {
    info: 'bg-cyan-500',
    warn: 'bg-yellow-500',
    error: 'bg-red-500',
    success: 'bg-emerald-500',
  }

  return (
    <div className="bg-[#161922] border border-[#1e2130] rounded-xl overflow-hidden">
      {/* Pipeline Visual Bar */}
      <div className="px-6 pt-5 pb-3 border-b border-[#1e2130]">
        <h3 className="text-sm font-semibold text-gray-200 mb-3">Pipeline</h3>
        <div className="flex items-center gap-1">
          {stages.map((stage, i) => (
            <div key={stage.id} className="flex items-center flex-1">
              <div className="flex flex-col items-center flex-1">
                <div
                  className={`w-full h-1.5 rounded-full transition-all ${
                    stage.status === 'success'
                      ? 'bg-emerald-500'
                      : stage.status === 'running'
                        ? 'bg-cyan-500 pipeline-active'
                        : stage.status === 'failed'
                          ? 'bg-red-500'
                          : 'bg-[#222639]'
                  }`}
                />
                <span
                  className={`text-[10px] mt-1.5 font-medium transition-colors ${
                    stage.status === 'success'
                      ? 'text-emerald-500'
                      : stage.status === 'running'
                        ? 'text-cyan-500'
                        : stage.status === 'failed'
                          ? 'text-red-500'
                          : 'text-gray-700'
                  }`}
                >
                  {stage.label}
                </span>
              </div>
              {i < stages.length - 1 && (
                <div className="h-1.5 w-6 shrink-0">
                  <svg
                    viewBox="0 0 24 8"
                    className="w-full h-full text-gray-800"
                    fill="currentColor"
                  >
                    <polygon points="0,0.5 23,3 0,7" />
                  </svg>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Log Terminal */}
      <div className="px-0">
        <div className="px-5 py-2.5 border-b border-[#1e2130] flex items-center justify-between">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
            Logs do Deploy
          </h3>
          <span className="text-[10px] text-gray-600 font-mono">
            {logs.length} etapas
          </span>
        </div>
        <div
          className="p-5 h-[220px] overflow-y-auto font-mono text-xs"
          style={{ background: '#0c0e14' }}
        >
          {logs.length === 0 ? (
            <div className="flex items-center justify-center h-full text-gray-700">
              Aguardando logs...
            </div>
          ) : (
            logs.map((log, i) => (
              <div key={i} className="flex items-start gap-2 leading-relaxed animate-fade-in py-0.5">
                <span className="text-gray-700 shrink-0 select-none">{log.timestamp}</span>
                <span
                  className={`shrink-0 w-1.5 h-1.5 rounded-full mt-1.5 ${levelDot[log.level]}`}
                />
                <span
                  className={
                    log.level === 'error'
                      ? 'text-red-400'
                      : log.level === 'warn'
                        ? 'text-yellow-400'
                        : log.level === 'success'
                          ? 'text-emerald-400'
                          : 'text-gray-300'
                  }
                >
                  {log.message}
                </span>
              </div>
            ))
          )}
          {deploy.status === 'running' && (
            <div className="flex items-center gap-2 mt-1">
              <span className="text-gray-700 animate-pulse">▊</span>
            </div>
          )}
          <div ref={logEndRef} />
        </div>
      </div>
    </div>
  )
}
