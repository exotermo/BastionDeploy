import { ExternalLink, GitBranch, User, Clock, ChevronDown, Maximize2, X } from 'lucide-react'
import { useDeploys, Deploy } from '../hooks/useDeploys'
import { useState } from 'react'
import { DeployLog } from './DeployLog'

type UIStatus = 'SUCCESS' | 'FAILED' | 'PENDING' | 'RUNNING'

const toUIStatus = (s: Deploy['status']): UIStatus =>
  s.toUpperCase() as UIStatus

const statusStyle: Record<UIStatus, string> = {
  SUCCESS: 'bg-emerald-500/10 text-emerald-400',
  FAILED: 'bg-red-500/10 text-red-400',
  PENDING: 'bg-amber-500/10 text-amber-400',
  RUNNING: 'bg-cyan-500/10 text-cyan-400',
}

const statusDot: Record<UIStatus, string> = {
  SUCCESS: 'bg-emerald-500',
  FAILED: 'bg-red-500',
  PENDING: 'bg-amber-500',
  RUNNING: 'bg-cyan-500 animate-pulse',
}

const formatTime = (dt: string) => {
  const diff = Math.floor((Date.now() - new Date(dt).getTime()) / 1000)
  if (diff < 60) return `${diff}s atrás`
  if (diff < 3600) return `${Math.floor(diff / 60)} min atrás`
  return `${Math.floor(diff / 3600)}h atrás`
}

function DeployDetailModal({
  deploy,
  onClose,
}: {
  deploy: Deploy
  onClose: () => void
}) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-8">
      <div className="absolute inset-0 bg-black/70" onClick={onClose} />
      <div className="relative w-full max-w-3xl bg-[#161922] border border-[#1e2130] rounded-xl overflow-hidden animate-fade-in">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-[#1e2130]">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-cyan-500/10 rounded-lg flex items-center justify-center">
              <GitBranch className="w-4 h-4 text-cyan-500" />
            </div>
            <div>
              <h3 className="text-sm font-semibold text-white">{deploy.app_name}</h3>
              <p className="text-[11px] text-gray-500">
                {deploy.id} · {deploy.triggered_by}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <span className="px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider">
              <span className="flex items-center gap-1.5">
                <span
                  className={`w-1.5 h-1.5 rounded-full ${statusDot[toUIStatus(deploy.status)]}`}
                />
                <span className={statusStyle[toUIStatus(deploy.status)]}>{toUIStatus(deploy.status)}</span>
              </span>
            </span>
            <button onClick={onClose} className="p-1.5 rounded-lg hover:bg-[#1e2130] text-gray-500 transition">
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Details bar */}
        <div className="px-6 py-3 bg-[#111319] border-b border-[#1e2130] flex items-center gap-6 text-[11px] text-gray-500">
          <span>
            Branch:{' '}
            <code className="text-gray-400 font-mono">
              {deploy.branch.replace('refs/heads/', '')}
            </code>
          </span>
          <span>
            Commit:{' '}
            <code className="text-gray-400 font-mono">
              {deploy.commit_sha.slice(0, 7)}
            </code>
          </span>
          <span className="flex items-center gap-1">
            <Clock className="w-3 h-3" />
            {formatTime(deploy.created_at)}
          </span>
        </div>

        {/* Pipeline & Logs */}
        <DeployLog deploy={deploy} />
      </div>
    </div>
  )
}

export function DeploymentsTable() {
  const { deploys, loading } = useDeploys()
  const [selectedDeploy, setSelectedDeploy] = useState<Deploy | null>(null)

  return (
    <div className="bg-[#161922] border border-[#1e2130] rounded-xl overflow-hidden">
      <div className="px-6 py-5 border-b border-[#1e2130] flex items-center justify-between">
        <div>
          <h2 className="text-base font-semibold text-white">Deploys Recentes</h2>
          <p className="text-[11px] text-gray-500 mt-0.5">
            Clique em um deploy para ver o pipeline detalhado
          </p>
        </div>
        <span className="text-[11px] text-gray-600 bg-[#1e2130] px-2.5 py-1 rounded-md">
          {deploys.length} deploys
        </span>
      </div>

      {loading ? (
        <div className="p-10 text-center text-gray-600 text-sm">Carregando deploys...</div>
      ) : deploys.length === 0 ? (
        <div className="p-10 text-center text-gray-600 text-sm">
          Nenhum deploy encontrado. Faça um push para o GitHub para iniciar.
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-[#1e2130]">
                {['App', 'Branch', 'Commit', 'Status', 'Autor', 'Tempo', ''].map((h) => (
                  <th
                    key={h}
                    className="text-left px-5 py-3 text-[10px] font-semibold text-gray-600 uppercase tracking-wider"
                  >
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {deploys.map((d) => {
                const uiStatus = toUIStatus(d.status)
                return (
                  <tr
                    key={d.id}
                    onClick={() => setSelectedDeploy(d)}
                    className="border-b border-[#1e2130]/50 hover:bg-[#1a1c27] transition-colors cursor-pointer"
                  >
                    <td className="px-5 py-3">
                      <span className="text-sm text-gray-200 font-medium">{d.app_name}</span>
                    </td>
                    <td className="px-5 py-3">
                      <code className="text-[12px] text-gray-500 font-mono">
                        {d.branch.replace('refs/heads/', '')}
                      </code>
                    </td>
                    <td className="px-5 py-3">
                      <code className="text-[12px] text-gray-500 bg-[#1a1c27] px-2 py-0.5 rounded font-mono">
                        {d.commit_sha.slice(0, 7)}
                      </code>
                    </td>
                    <td className="px-5 py-3">
                      <span className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[11px] font-semibold">
                        <span
                          className={`w-1.5 h-1.5 rounded-full ${statusDot[uiStatus]}`}
                        />
                        {uiStatus}
                      </span>
                    </td>
                    <td className="px-5 py-3">
                      <div className="flex items-center gap-1.5 text-[13px] text-gray-500">
                        <User className="w-3.5 h-3.5" />
                        {d.triggered_by}
                      </div>
                    </td>
                    <td className="px-5 py-3 text-[12px] text-gray-600">
                      {formatTime(d.created_at)}
                    </td>
                    <td className="px-5 py-3">
                      <Maximize2 className="w-3.5 h-3.5 text-gray-700" />
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      {selectedDeploy && (
        <DeployDetailModal
          deploy={selectedDeploy}
          onClose={() => setSelectedDeploy(null)}
        />
      )}
    </div>
  )
}
