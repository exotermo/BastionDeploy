import { useState } from 'react'
import { useDeploys, Deploy } from '../hooks/useDeploys'
import { DeployLog } from '../components/DeployLog'
import { GitBranch, User, Clock, ArrowLeft, Maximize2, X } from 'lucide-react'

type UIStatus = 'SUCCESS' | 'FAILED' | 'PENDING' | 'RUNNING'
const toUIStatus = (s: Deploy['status']): UIStatus => s.toUpperCase() as UIStatus

const statusDot: Record<UIStatus, string> = {
  SUCCESS: 'bg-emerald-500',
  FAILED: 'bg-red-500',
  PENDING: 'bg-amber-500',
  RUNNING: 'bg-cyan-500 animate-pulse',
}

const statusText: Record<UIStatus, string> = {
  SUCCESS: 'text-emerald-500',
  FAILED: 'text-red-500',
  PENDING: 'text-amber-500',
  RUNNING: 'text-cyan-500',
}

const formatTime = (dt: string) => {
  const diff = Math.floor((Date.now() - new Date(dt).getTime()) / 1000)
  if (diff < 60) return `${diff}s atrás`
  if (diff < 3600) return `${Math.floor(diff / 60)} min atrás`
  return `${Math.floor(diff / 3600)}h atrás`
}

export function DeployPage() {
  const { deploys, loading } = useDeploys()
  const [selected, setSelected] = useState<Deploy | null>(null)

  return (
    <>
      <div className="mb-6">
        <h1 className="text-xl font-bold text-white">Histórico de Deploys</h1>
        <p className="text-[13px] text-gray-500 mt-0.5">
          Clique em um deploy para ver o pipeline completo
        </p>
      </div>

      {/* Deploy list */}
      {loading ? (
        <div className="text-center py-12 text-gray-600">Carregando...</div>
      ) : deploys.length === 0 ? (
        <div className="bg-[#161922] border border-[#1e2130] rounded-xl p-12 text-center">
          <p className="text-gray-500 mb-2">Nenhum deploy encontrado</p>
          <p className="text-sm text-gray-700">Faça um push para o GitHub para começar.</p>
        </div>
      ) : (
        <div className="space-y-2">
          {deploys.map((d) => {
            const s = toUIStatus(d.status)
            return (
              <div
                key={d.id}
                onClick={() => setSelected(d)}
                className="bg-[#161922] border border-[#1e2130] rounded-xl p-4 hover:border-[#2a2d42] cursor-pointer transition-all flex items-center justify-between"
              >
                <div className="flex items-center gap-4 min-w-0">
                  <div className={`w-2 h-2 rounded-full shrink-0 ${statusDot[s]}`} />
                  <div className="min-w-0">
                    <p className="text-sm text-gray-200 font-medium truncate">{d.app_name}</p>
                    <p className="text-[11px] text-gray-600 flex items-center gap-1.5 mt-0.5">
                      <code className="font-mono">{d.branch.replace('refs/heads/', '')}</code>
                      <span className="font-mono bg-[#1a1c27] px-1.5 rounded">{d.commit_sha.slice(0, 7)}</span>
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-5 shrink-0">
                  <div className="flex items-center gap-1.5 text-[12px] text-gray-600">
                    <User className="w-3.5 h-3.5" />{d.triggered_by}
                  </div>
                  <div className="flex items-center gap-1 text-[12px] text-gray-600">
                    <Clock className="w-3.5 h-3.5" />{formatTime(d.created_at)}
                  </div>
                  <span className={`text-[11px] font-semibold ${statusText[s]}`}>{s}</span>
                  <Maximize2 className="w-4 h-4 text-gray-700" />
                </div>
              </div>
            )
          })}
        </div>
      )}

      {/* Detail modal */}
      {selected && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-6">
          <div className="absolute inset-0 bg-black/70" onClick={() => setSelected(null)} />
          <div className="relative w-full max-w-3xl bg-[#161922] border border-[#1e2130] rounded-xl overflow-hidden">
            <div className="flex items-center justify-between px-5 py-4 border-b border-[#1e2130]">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-cyan-500/10 rounded-lg flex items-center justify-center">
                  <GitBranch className="w-4 h-4 text-cyan-500" />
                </div>
                <div>
                  <h3 className="text-sm font-semibold text-white">{selected.app_name}</h3>
                  <p className="text-[11px] text-gray-600">{selected.triggered_by}</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <span className={`text-[11px] font-semibold ${statusText[toUIStatus(selected.status)]}`}>
                  {toUIStatus(selected.status)}
                </span>
                <button onClick={() => setSelected(null)} className="p-1 rounded hover:bg-[#1e2130] text-gray-500">
                  <X className="w-4 h-4" />
                </button>
              </div>
            </div>
            <div className="px-4 py-2 bg-[#111319] border-b border-[#1e2130] flex items-center gap-5 text-[11px] text-gray-500">
              <span>Branch: <code className="text-gray-400 font-mono">{selected.branch.replace('refs/heads/', '')}</code></span>
              <span>Commit: <code className="text-gray-400 font-mono">{selected.commit_sha.slice(0, 7)}</code></span>
              <span className="flex items-center gap-1"><Clock className="w-3 h-3" />{formatTime(selected.created_at)}</span>
            </div>
            <DeployLog deploy={selected} />
          </div>
        </div>
      )}
    </>
  )
}
